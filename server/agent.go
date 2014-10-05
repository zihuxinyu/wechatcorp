// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"github.com/chanxuehong/wechatcorp/message/passive/request"
	"net/http"
)

// 企业号应用对外暴露的接口
type Agent interface {
	GetCorpId() string // 获取应用所属的企业号Id
	GetAgentId() int64 // 获取应用的Id
	GetToken() string  // 对应后台的设置的 Token
	GetAESKey() []byte // 32 bytes 的 AES 加密 Key

	// 未知类型的消息处理方法
	//  rawXMLMsg 是解密后的明文 xml 消息体
	//  timestamp 是请求 URL 中的时间戳
	//  nonce     是请求 URL 中的随机数
	//  random    是请求 http body 中的密文消息加密时所用的 random
	//  r *http.Request 的 Body 已经读取过了, 不要再读取了, 但是可以获取其他信息, 比如 user-agent
	ServeUnknownMsg(w http.ResponseWriter, r *http.Request, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)

	// 消息处理函数
	//  rawXMLMsg 是解密后的明文 xml 消息体
	//  timestamp 是请求 URL 中的时间戳
	//  nonce     是请求 URL 中的随机数
	//  random    是请求 http body 中的密文消息加密时所用的 random
	//  r *http.Request 的 Body 已经读取过了, 不要再读取了, 但是可以获取其他信息, 比如 user-agent
	ServeTextMsg(w http.ResponseWriter, r *http.Request, msg *request.Text, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ServeImageMsg(w http.ResponseWriter, r *http.Request, msg *request.Image, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ServeVoiceMsg(w http.ResponseWriter, r *http.Request, msg *request.Voice, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ServeVideoMsg(w http.ResponseWriter, r *http.Request, msg *request.Video, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ServeLocationMsg(w http.ResponseWriter, r *http.Request, msg *request.Location, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)

	// 事件处理函数
	//  rawXMLMsg 是解密后的明文 xml 消息体
	//  timestamp 是请求 URL 中的时间戳
	//  nonce     是请求 URL 中的随机数
	//  random    是请求 http body 中的密文消息加密时所用的 random
	//  r *http.Request 的 Body 已经读取过了, 不要再读取了, 但是可以获取其他信息, 比如 user-agent
	ServeSubscribeEvent(w http.ResponseWriter, r *http.Request, event *request.SubscribeEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ServeUnsubscribeEvent(w http.ResponseWriter, r *http.Request, event *request.UnsubscribeEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ServeLocationEvent(w http.ResponseWriter, r *http.Request, event *request.LocationEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ServeMenuClickEvent(w http.ResponseWriter, r *http.Request, event *request.MenuClickEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ServeMenuViewEvent(w http.ResponseWriter, r *http.Request, event *request.MenuViewEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
}
