// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/chanxuehong/wechatcorp/message/passive/request"
	"net/http"
	"net/url"
	"strconv"
)

// net/http.Handler 的实现.
//  NOTE: 只能处理一个 Agent 的消息
type SingleHttpHandler struct {
	invalidRequestHandler InvalidRequestHandler
	agentMsgHandler       AgentMsgHandler
}

func NewSingleHttpHandler(invalidRequestHandler InvalidRequestHandler,
	agentMsgHandler AgentMsgHandler) *SingleHttpHandler {

	if agentMsgHandler == nil {
		panic("agentMsgHandler == nil")
	}
	if invalidRequestHandler == nil {
		invalidRequestHandler = InvalidRequestHandlerFunc(DefaultInvalidRequestHandlerFunc)
	}

	return &SingleHttpHandler{
		invalidRequestHandler: invalidRequestHandler,
		agentMsgHandler:       agentMsgHandler,
	}
}

func (this *SingleHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	invalidRequestHandler := this.invalidRequestHandler
	agentMsgHandler := this.agentMsgHandler

	switch r.Method {
	case "POST": // 处理从微信服务器推送过来的消息(事件) ==============================
		var requestHttpBody request.RequestHttpBody
		if err := xml.NewDecoder(r.Body).Decode(&requestHttpBody); err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		if r.URL == nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("r.URL == nil"))
			return
		}

		urlValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		signature := urlValues.Get("msg_signature")
		if signature == "" {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("msg_signature is empty"))
			return
		}

		const signatureLen = sha1.Size * 2
		if len(signature) != signatureLen {
			invalidRequestHandler.ServeInvalidRequest(w, r, fmt.Errorf("the length of signature mismatch, have: %d, want: %d", len(signature), signatureLen))
			return
		}

		timestampStr := urlValues.Get("timestamp")
		if timestampStr == "" {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("timestamp is empty"))
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, fmt.Errorf("can not parse timestamp(==%q) to int64, error: %s", timestampStr, err.Error()))
			return
		}

		nonce := urlValues.Get("nonce")
		if nonce == "" {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("nonce is empty"))
			return
		}

		signaturex := agentMsgHandler.Signature(timestampStr, nonce, requestHttpBody.EncryptMsg)
		if subtle.ConstantTimeCompare([]byte(signature), []byte(signaturex)) != 1 {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("check signature failed"))
			return
		}

		EncryptMsgBytes, err := base64.StdEncoding.DecodeString(requestHttpBody.EncryptMsg)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		random, rawXMLMsg, err := agentMsgHandler.DecryptMsg(EncryptMsgBytes)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		var msgReq request.Request
		if err := xml.Unmarshal(rawXMLMsg, &msgReq); err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// request router, 可一个根据自己的实际业务调整顺序!
		switch msgReq.MsgType {
		case request.MSG_TYPE_TEXT:
			agentMsgHandler.TextMsgHandler(w, r, msgReq.Text(), rawXMLMsg, timestamp, nonce, random)

		case request.MSG_TYPE_EVENT:
			// event router
			switch msgReq.Event {
			case request.EVENT_TYPE_LOCATION:
				agentMsgHandler.LocationEventHandler(w, r, msgReq.LocationEvent(), rawXMLMsg, timestamp, nonce, random)

			case request.EVENT_TYPE_CLICK:
				agentMsgHandler.MenuClickEventHandler(w, r, msgReq.MenuClickEvent(), rawXMLMsg, timestamp, nonce, random)

			case request.EVENT_TYPE_VIEW:
				agentMsgHandler.MenuViewEventHandler(w, r, msgReq.MenuViewEvent(), rawXMLMsg, timestamp, nonce, random)

			case request.EVENT_TYPE_SUBSCRIBE:
				agentMsgHandler.SubscribeEventHandler(w, r, msgReq.SubscribeEvent(), rawXMLMsg, timestamp, nonce, random)

			case request.EVENT_TYPE_UNSUBSCRIBE:
				agentMsgHandler.UnsubscribeEventHandler(w, r, msgReq.UnsubscribeEvent(), rawXMLMsg, timestamp, nonce, random)

			default: // unknown event
				agentMsgHandler.UnknownMsgHandler(w, r, rawXMLMsg, timestamp, nonce, random)
			}

		case request.MSG_TYPE_VOICE:
			agentMsgHandler.VoiceMsgHandler(w, r, msgReq.Voice(), rawXMLMsg, timestamp, nonce, random)

		case request.MSG_TYPE_LOCATION:
			agentMsgHandler.LocationMsgHandler(w, r, msgReq.Location(), rawXMLMsg, timestamp, nonce, random)

		case request.MSG_TYPE_IMAGE:
			agentMsgHandler.ImageMsgHandler(w, r, msgReq.Image(), rawXMLMsg, timestamp, nonce, random)

		case request.MSG_TYPE_VIDEO:
			agentMsgHandler.VideoMsgHandler(w, r, msgReq.Video(), rawXMLMsg, timestamp, nonce, random)

		default: // unknown request message type
			agentMsgHandler.UnknownMsgHandler(w, r, rawXMLMsg, timestamp, nonce, random)
		}

	case "GET": // 首次验证 ======================================================
		if r.URL == nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("r.URL == nil"))
			return
		}

		urlValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		signature := urlValues.Get("msg_signature")
		if signature == "" {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("msg_signature is empty"))
			return
		}

		const signatureLen = sha1.Size * 2
		if len(signature) != signatureLen {
			invalidRequestHandler.ServeInvalidRequest(w, r, fmt.Errorf("the length of signature mismatch, have: %d, want: %d", len(signature), signatureLen))
			return
		}

		timestamp := urlValues.Get("timestamp")
		if timestamp == "" {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("timestamp is empty"))
			return
		}

		nonce := urlValues.Get("nonce")
		if nonce == "" {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("nonce is empty"))
			return
		}

		EncryptMsg := urlValues.Get("echostr")
		if EncryptMsg == "" {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("echostr is empty"))
			return
		}

		EncryptMsgBytes, err := base64.StdEncoding.DecodeString(EncryptMsg)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		signaturex := agentMsgHandler.Signature(timestamp, nonce, EncryptMsg)
		if subtle.ConstantTimeCompare([]byte(signature), []byte(signaturex)) != 1 {
			invalidRequestHandler.ServeInvalidRequest(w, r, errors.New("check signature failed"))
			return
		}

		_, echostr, err := agentMsgHandler.DecryptMsg(EncryptMsgBytes)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		w.Write(echostr)
	}
}
