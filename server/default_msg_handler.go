// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"github.com/chanxuehong/wechatcorp/message/passive/request"
	"net/http"
)

var _ AgentMsgHandler = new(DefaultAgentMsgHandler)

type DefaultAgentMsgHandler struct {
	CorpId   string
	AgentId  string
	Token    string
	AESKey   []byte
	CipherIV []byte
}

func (handler *DefaultAgentMsgHandler) Init(CorpId, AgentId, Token string, AESKey []byte) {
	if len(AESKey) != 32 {
		panic("the length of AESKey must be equal 32")
	}

	handler.CorpId = CorpId
	handler.AgentId = AgentId
	handler.Token = Token
	handler.AESKey = AESKey
	handler.CipherIV = AESKey[:16]
}

func (handler *DefaultAgentMsgHandler) UnknownMsgHandler(w http.ResponseWriter, r *http.Request, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) TextMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Text, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) ImageMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Image, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) VoiceMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Voice, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) VideoMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Video, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) LocationMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Location, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) SubscribeEventHandler(w http.ResponseWriter, r *http.Request, event *request.SubscribeEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) UnsubscribeEventHandler(w http.ResponseWriter, r *http.Request, event *request.UnsubscribeEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) LocationEventHandler(w http.ResponseWriter, r *http.Request, event *request.LocationEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) MenuClickEventHandler(w http.ResponseWriter, r *http.Request, event *request.MenuClickEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
func (handler *DefaultAgentMsgHandler) MenuViewEventHandler(w http.ResponseWriter, r *http.Request, event *request.MenuViewEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
}
