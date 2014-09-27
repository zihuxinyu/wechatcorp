// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"github.com/chanxuehong/wechatcorp/message/passive/response"
	"io"
	"strconv"
)

// 把 text 回复消息 msg 写入 writer w
func (handler *DefaultAgentMsgHandler) WriteText(w io.Writer, msg *response.Text, timestamp, nonce int64, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	return handler.writeResponse(w, msg, timestamp, nonce, random)
}

// 把 image 回复消息 msg 写入 writer w
func (handler *DefaultAgentMsgHandler) WriteImage(w io.Writer, msg *response.Image, timestamp, nonce int64, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	return handler.writeResponse(w, msg, timestamp, nonce, random)
}

// 把 voice 回复消息 msg 写入 writer w
func (handler *DefaultAgentMsgHandler) WriteVoice(w io.Writer, msg *response.Voice, timestamp, nonce int64, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	return handler.writeResponse(w, msg, timestamp, nonce, random)
}

// 把 video 回复消息 msg 写入 writer w
func (handler *DefaultAgentMsgHandler) WriteVideo(w io.Writer, msg *response.Video, timestamp, nonce int64, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	return handler.writeResponse(w, msg, timestamp, nonce, random)
}

// 把 news 回复消息 msg 写入 writer w
func (handler *DefaultAgentMsgHandler) WriteNews(w io.Writer, msg *response.News, timestamp, nonce int64, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	if err = msg.CheckValid(); err != nil {
		return
	}
	return handler.writeResponse(w, msg, timestamp, nonce, random)
}

func (handler *DefaultAgentMsgHandler) writeResponse(w io.Writer, msg interface{}, timestamp, nonce int64, random []byte) (err error) {
	msgBytes, err := xml.Marshal(msg)
	if err != nil {
		return
	}

	EncryptMsg := handler.EncryptMsg(random, msgBytes)
	base64EncryptMsg := base64.StdEncoding.EncodeToString(EncryptMsg)

	var ResponseHttpBody response.ResponseHttpBody
	ResponseHttpBody.EncryptMsg = base64EncryptMsg
	ResponseHttpBody.TimeStamp = timestamp
	ResponseHttpBody.Nonce = nonce

	timestampStr := strconv.FormatInt(timestamp, 10)
	nonceStr := strconv.FormatInt(nonce, 10)
	ResponseHttpBody.Signature = string(handler.Signature(timestampStr, nonceStr, base64EncryptMsg))

	return xml.NewEncoder(w).Encode(&ResponseHttpBody)
}
