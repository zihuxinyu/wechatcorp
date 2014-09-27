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
	"github.com/chanxuehong/wechatcorp/message/passive/request"
	"net/http"
	"net/url"
	"strconv"
)

// net/http.Handler 的实现
type AgentHttpHandler struct {
	AgentMsgHandler AgentMsgHandler
}

func (handler AgentHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST": // 处理从微信服务器推送过来的消息(事件) ==============================
		var RequestHttpBody request.RequestHttpBody
		if err := xml.NewDecoder(r.Body).Decode(&RequestHttpBody); err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		if RequestHttpBody.CorpID != handler.AgentMsgHandler.GetCorpId() {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("CorpId mismatch"))
			return
		}
		if RequestHttpBody.AgentId != handler.AgentMsgHandler.GetAgentId() {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("AgentId mismatch"))
			return
		}

		if r.URL == nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("r.URL == nil"))
			return
		}

		urlValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		signature := urlValues.Get("msg_signature")
		if signature == "" {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("msg_signature is empty"))
			return
		}

		const signatureLen = sha1.Size * 2
		if len(signature) != signatureLen {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("check signature failed"))
			return
		}

		timestampStr := urlValues.Get("timestamp")
		if timestampStr == "" {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("timestamp is empty"))
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		nonceStr := urlValues.Get("nonce")
		if nonceStr == "" {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("nonce is empty"))
			return
		}

		nonce, err := strconv.ParseInt(nonceStr, 10, 64)
		if err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		signaturex := handler.AgentMsgHandler.Signature(timestampStr, nonceStr, RequestHttpBody.EncryptMsg)
		// 采用 subtle.ConstantTimeCompare 是防止 计时攻击!
		if subtle.ConstantTimeCompare([]byte(signature), []byte(signaturex)) != 1 {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("check signature failed"))
			return
		}

		EncryptMsgBytes, err := base64.StdEncoding.DecodeString(RequestHttpBody.EncryptMsg)
		if err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		random, rawXMLMsg, err := handler.AgentMsgHandler.DecryptMsg(EncryptMsgBytes)
		if err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		var msgReq request.Request
		if err := xml.Unmarshal(rawXMLMsg, &msgReq); err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		// request router, 可一个根据自己的实际业务调整顺序!
		switch msgReq.MsgType {
		case request.MSG_TYPE_TEXT:
			handler.AgentMsgHandler.TextMsgHandler(w, r, msgReq.Text(), rawXMLMsg, timestamp, nonce, random)

		case request.MSG_TYPE_EVENT:
			// event router
			switch msgReq.Event {
			case request.EVENT_TYPE_LOCATION:
				handler.AgentMsgHandler.LocationEventHandler(w, r, msgReq.LocationEvent(), rawXMLMsg, timestamp, nonce, random)

			case request.EVENT_TYPE_CLICK:
				handler.AgentMsgHandler.MenuClickEventHandler(w, r, msgReq.MenuClickEvent(), rawXMLMsg, timestamp, nonce, random)

			case request.EVENT_TYPE_VIEW:
				handler.AgentMsgHandler.MenuViewEventHandler(w, r, msgReq.MenuViewEvent(), rawXMLMsg, timestamp, nonce, random)

			case request.EVENT_TYPE_SUBSCRIBE:
				handler.AgentMsgHandler.SubscribeEventHandler(w, r, msgReq.SubscribeEvent(), rawXMLMsg, timestamp, nonce, random)

			case request.EVENT_TYPE_UNSUBSCRIBE:
				handler.AgentMsgHandler.UnsubscribeEventHandler(w, r, msgReq.UnsubscribeEvent(), rawXMLMsg, timestamp, nonce, random)

			default: // unknown event
				handler.AgentMsgHandler.UnknownMsgHandler(w, r, rawXMLMsg, timestamp, nonce, random)
			}

		case request.MSG_TYPE_VOICE:
			handler.AgentMsgHandler.VoiceMsgHandler(w, r, msgReq.Voice(), rawXMLMsg, timestamp, nonce, random)

		case request.MSG_TYPE_LOCATION:
			handler.AgentMsgHandler.LocationMsgHandler(w, r, msgReq.Location(), rawXMLMsg, timestamp, nonce, random)

		case request.MSG_TYPE_IMAGE:
			handler.AgentMsgHandler.ImageMsgHandler(w, r, msgReq.Image(), rawXMLMsg, timestamp, nonce, random)

		case request.MSG_TYPE_VIDEO:
			handler.AgentMsgHandler.VideoMsgHandler(w, r, msgReq.Video(), rawXMLMsg, timestamp, nonce, random)

		default: // unknown request message type
			handler.AgentMsgHandler.UnknownMsgHandler(w, r, rawXMLMsg, timestamp, nonce, random)
		}

	case "GET": // 首次验证 ======================================================
		if r.URL == nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("r.URL == nil"))
			return
		}

		urlValues, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		signature := urlValues.Get("msg_signature")
		if signature == "" {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("msg_signature is empty"))
			return
		}

		const signatureLen = sha1.Size * 2
		if len(signature) != signatureLen {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("check signature failed"))
			return
		}

		timestamp := urlValues.Get("timestamp")
		if timestamp == "" {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("timestamp is empty"))
			return
		}

		nonce := urlValues.Get("nonce")
		if nonce == "" {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("nonce is empty"))
			return
		}

		EncryptMsg := urlValues.Get("echostr")
		if EncryptMsg == "" {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("echostr is empty"))
			return
		}

		signaturex := handler.AgentMsgHandler.Signature(timestamp, nonce, EncryptMsg)
		// 采用 subtle.ConstantTimeCompare 是防止 计时攻击!
		if subtle.ConstantTimeCompare([]byte(signature), []byte(signaturex)) != 1 {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, errors.New("check signature failed"))
			return
		}

		EncryptMsgBytes, err := base64.StdEncoding.DecodeString(EncryptMsg)
		if err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		_, echostr, err := handler.AgentMsgHandler.DecryptMsg(EncryptMsgBytes)
		if err != nil {
			handler.AgentMsgHandler.InvalidRequestHandler(w, r, err)
			return
		}

		w.Write(echostr)
	}
}
