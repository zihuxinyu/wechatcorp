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
)

// Agent 的前端, 负责处理 http 请求, net/http.Handler 的实现
//  NOTE: 只能处理一个企业号应用的消息
type AgentFrontend struct {
	agent                 Agent
	invalidRequestHandler InvalidRequestHandler
}

func NewAgentFrontend(agent Agent, invalidRequestHandler InvalidRequestHandler) *AgentFrontend {
	if agent == nil {
		panic("agent == nil")
	}
	if invalidRequestHandler == nil {
		invalidRequestHandler = InvalidRequestHandlerFunc(defaultInvalidRequestHandlerFunc)
	}

	return &AgentFrontend{
		agent: agent,
		invalidRequestHandler: invalidRequestHandler,
	}
}

func (this *AgentFrontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	agent := this.agent
	invalidRequestHandler := this.invalidRequestHandler

	switch r.Method {
	case "POST": // 处理从微信服务器推送过来的消息(事件) ==============================
		signature1, timestampStr, nonce, err := parsePostURLQuery(r.URL)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, fmt.Errorf("can not parse timestamp(==%q) to int64, error: %s", timestampStr, err.Error()))
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

		wantCorpId := agent.GetCorpId()
		if subtle.ConstantTimeCompare([]byte(requestHttpBody.CorpId), []byte(wantCorpId)) != 1 {
			err = fmt.Errorf("the message RequestHttpBody's ToUserName mismatch, have: %s, want: %s", requestHttpBody.CorpId, wantCorpId)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}
		wantAgentId := agent.GetAgentId()
		if subtle.ConstantTimeCompare([]byte(requestHttpBody.AgentId), []byte(wantAgentId)) != 1 {
			err = fmt.Errorf("the message RequestHttpBody's AgentId mismatch, have: %s, want: %s", requestHttpBody.AgentId, wantAgentId)
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

		msgDispatch(w, r, &msgReq, rawXMLMsg, timestamp, nonce, random, agent)

	case "GET": // 首次验证 ======================================================
		signature1, timestamp, nonce, encryptMsg, err := parseGetURLQuery(r.URL)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
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
