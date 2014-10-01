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
func (this *DefaultAgent) WriteText(w io.Writer, msg *response.Text, timestamp int64, nonce string, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	return this.writeResponse(w, msg, timestamp, nonce, random)
}

// 把 image 回复消息 msg 写入 writer w
func (this *DefaultAgent) WriteImage(w io.Writer, msg *response.Image, timestamp int64, nonce string, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	return this.writeResponse(w, msg, timestamp, nonce, random)
}

// 把 voice 回复消息 msg 写入 writer w
func (this *DefaultAgent) WriteVoice(w io.Writer, msg *response.Voice, timestamp int64, nonce string, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	return this.writeResponse(w, msg, timestamp, nonce, random)
}

// 把 video 回复消息 msg 写入 writer w
func (this *DefaultAgent) WriteVideo(w io.Writer, msg *response.Video, timestamp int64, nonce string, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	return this.writeResponse(w, msg, timestamp, nonce, random)
}

// 把 news 回复消息 msg 写入 writer w
func (this *DefaultAgent) WriteNews(w io.Writer, msg *response.News, timestamp int64, nonce string, random []byte) (err error) {
	if w == nil {
		return errors.New("w == nil")
	}
	if msg == nil {
		return errors.New("msg == nil")
	}
	if err = msg.CheckValid(); err != nil {
		return
	}
	return this.writeResponse(w, msg, timestamp, nonce, random)
}

func (this *DefaultAgent) writeResponse(w io.Writer, msg interface{}, timestamp int64, nonce string, random []byte) (err error) {
	rawXMLMsg, err := xml.Marshal(msg)
	if err != nil {
		return
	}

	EncryptMsg := encryptMsg(random, rawXMLMsg, this.CorpId, this.AESKey)
	base64EncryptMsg := base64.StdEncoding.EncodeToString(EncryptMsg)

	var responseHttpBody response.ResponseHttpBody
	responseHttpBody.EncryptMsg = base64EncryptMsg
	responseHttpBody.TimeStamp = timestamp
	responseHttpBody.Nonce = nonce

	timestampStr := strconv.FormatInt(timestamp, 10)
	responseHttpBody.Signature = signature(this.Token, timestampStr, nonce, base64EncryptMsg)

	return xml.NewEncoder(w).Encode(&responseHttpBody)
}
