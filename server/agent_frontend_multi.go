// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/chanxuehong/wechatcorp/message/passive/request"
	"net/http"
	"strconv"
	"sync"
)

// 定义回调 URL 上指定 Agent 的查询参数名
const URLQueryAgentKeyName = "agentkey"

// 多个 Agent 的前端, 负责处理 http 请求, net/http.Handler 的实现
//
//  NOTE:
//  MultiAgentFrontend 可以处理多个公众号的消息（事件），但是要求在回调 URL 上加上一个查询
//  参数，一般为 agentkey（参考常量 URLQueryAgentKeyName），这个参数的值就是 MultiAgentFrontend
//  索引 Agent 的 key。
//  例如回调 URL 为 http://www.xxx.com/weixin?agentkey=1234567890，那么就可以在后端调用
//
//    MultiAgentFrontend.SetAgent("1234567890", agent)
//
//  来增加一个 Agent 来处理 agentkey=1234567890 的消息（事件）。
//
//  MultiAgentFrontend 并发安全，可以在运行中动态增加和删除 Agent。
type MultiAgentFrontend struct {
	rwmutex               sync.RWMutex
	agentMap              map[string]Agent
	invalidRequestHandler InvalidRequestHandler
}

// 设置 InvalidRequestHandler, 如果 handler == nil 则使用默认的 DefaultInvalidRequestHandlerFunc
func (this *MultiAgentFrontend) SetInvalidRequestHandler(handler InvalidRequestHandler) {
	this.rwmutex.Lock()
	if handler == nil {
		this.invalidRequestHandler = InvalidRequestHandlerFunc(defaultInvalidRequestHandlerFunc)
	} else {
		this.invalidRequestHandler = handler
	}
	this.rwmutex.Unlock()
}

// 添加（设置） agentkey-agent pair, 如果 agent == nil 则不做任何操作
func (this *MultiAgentFrontend) SetAgent(agentkey string, agent Agent) {
	if agent == nil {
		return
	}

	this.rwmutex.Lock()
	if this.agentMap == nil {
		this.agentMap = make(map[string]Agent)
	}
	this.agentMap[agentkey] = agent
	this.rwmutex.Unlock()
}

// 删除 agentkey 对应的 Agent
func (this *MultiAgentFrontend) DeleteAgent(agentkey string) {
	this.rwmutex.Lock()
	delete(this.agentMap, agentkey)
	this.rwmutex.Unlock()
}

// 删除所有的 Agent
func (this *MultiAgentFrontend) DeleteAllAgent() {
	this.rwmutex.Lock()
	this.agentMap = nil
	this.rwmutex.Unlock()
}

func (this *MultiAgentFrontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.rwmutex.RLock()
	defer this.rwmutex.RUnlock()

	invalidRequestHandler := this.invalidRequestHandler
	if invalidRequestHandler == nil {
		invalidRequestHandler = InvalidRequestHandlerFunc(defaultInvalidRequestHandlerFunc)
	}
	if len(this.agentMap) == 0 {
		invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("no Agent"))
		return
	}

	switch r.Method {
	case "POST": // 处理从微信服务器推送过来的消息(事件) ==============================
		agentkey, signature1, timestampStr, nonce, err := parsePostURLQueryEx(r.URL)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, fmt.Errorf("can not parse timestamp(==%q) to int64, error: %s", timestampStr, err.Error()))
			return
		}

		agent := this.agentMap[agentkey]
		if agent == nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, fmt.Errorf("Not found Agent for %s == %s", URLQueryAgentKeyName, agentkey))
			return
		}

		var requestHttpBody request.RequestHttpBody
		if err := xml.NewDecoder(r.Body).Decode(&requestHttpBody); err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		signature2 := signature(agent.GetToken(), timestampStr, nonce, requestHttpBody.EncryptMsg)
		if subtle.ConstantTimeCompare([]byte(signature1), []byte(signature2)) != 1 {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("check signature failed"))
			return
		}

		if wantCorpId := agent.GetCorpId(); subtle.ConstantTimeCompare([]byte(requestHttpBody.CorpId), []byte(wantCorpId)) != 1 {
			err = fmt.Errorf("the message RequestHttpBody's ToUserName mismatch, have: %s, want: %s", requestHttpBody.CorpId, wantCorpId)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		wantAgentId := agent.GetAgentId()

		if requestHttpBody.AgentId != wantAgentId && requestHttpBody.AgentId != 0 {
			err = fmt.Errorf("the message RequestHttpBody's AgentId mismatch, have: %d, want: %d", requestHttpBody.AgentId, wantAgentId)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		EncryptMsgBytes, err := base64.StdEncoding.DecodeString(requestHttpBody.EncryptMsg)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		random, rawXMLMsg, err := decryptMsg(EncryptMsgBytes, agent.GetCorpId(), agent.GetAESKey())
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		var msgReq request.Request
		if err := xml.Unmarshal(rawXMLMsg, &msgReq); err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		if requestHttpBody.CorpId != msgReq.ToUserName {
			err = fmt.Errorf("the RequestHttpBody's ToUserName(==%d) mismatch the Request's ToUserName(==%d)", requestHttpBody.CorpId, msgReq.ToUserName)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		if requestHttpBody.AgentId != msgReq.AgentId {
			err = fmt.Errorf("the RequestHttpBody's AgentId(==%d) mismatch the Request's AgengId(==%d)", requestHttpBody.AgentId, msgReq.AgentId)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// 此时要么 msgReq.AgentId == wantAgentId, 要么 msgReq.AgentId == 0

		if msgReq.AgentId == 0 {
			// 订阅/取消订阅 整个企业号
			if msgReq.MsgType == request.MSG_TYPE_EVENT &&
				(msgReq.Event == request.EVENT_TYPE_SUBSCRIBE || msgReq.Event == request.EVENT_TYPE_UNSUBSCRIBE) {
				// do nothing
			} else {
				err = fmt.Errorf("the message Request's AgentId mismatch, have: %d, want: %d", msgReq.AgentId, wantAgentId)
				invalidRequestHandler.ServeInvalidRequest(w, r, err)
				return
			}
		}

		msgDispatch(w, r, &msgReq, rawXMLMsg, timestamp, nonce, random, agent)

	case "GET": // 首次验证 ======================================================
		agentkey, signature1, timestamp, nonce, encryptMsg, err := parseGetURLQueryEx(r.URL)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		agent := this.agentMap[agentkey]
		if agent == nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, fmt.Errorf("Not found Agent for %s == %s", URLQueryAgentKeyName, agentkey))
			return
		}

		signature2 := signature(agent.GetToken(), timestamp, nonce, encryptMsg)
		if subtle.ConstantTimeCompare([]byte(signature1), []byte(signature2)) != 1 {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("check signature failed"))
			return
		}

		EncryptMsgBytes, err := base64.StdEncoding.DecodeString(encryptMsg)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		_, echostr, err := decryptMsg(EncryptMsgBytes, agent.GetCorpId(), agent.GetAESKey())
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		w.Write(echostr)
	}
}
