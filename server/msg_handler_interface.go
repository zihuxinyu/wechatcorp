// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"github.com/chanxuehong/wechatcorp/message/passive/request"
	"net/http"
)

// 应用的消息处理接口
type AgentMsgHandler interface {
	// 获取应用的企业号 id
	GetCorpId() string

	// 获取应用的 id
	GetAgentId() string

	// 生成签名
	Signature(timestamp, nonce, EncryptMsg string) (signature string)

	// 加密.
	// random 的长度为 16, 你也可以不使用参数指定的值 random, 可以自己生成!
	EncryptMsg(random, rawXMLMsg []byte) (EncryptMsg []byte)

	// 解密, 要验证 corp id 的正确性
	DecryptMsg(EncryptMsg []byte) (random, rawXMLMsg []byte, err error)

	// 非法的请求处理方法, err 是出错信息
	InvalidRequestHandler(w http.ResponseWriter, r *http.Request, err error)

	// 未知类型的消息处理方法
	//  rawXMLMsg 是解密后的明文 xml 消息体
	//  timestamp 是请求中的时间戳
	//  nonce     是请求中的随机数
	//  random    是请求中的消息体加密的 random
	UnknownMsgHandler(w http.ResponseWriter, r *http.Request, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)

	// 消息处理函数
	//  rawXMLMsg 是解密后的明文 xml 消息体
	//  timestamp 是请求中的时间戳
	//  nonce     是请求中的随机数
	//  random    是请求中的消息体加密的 random
	TextMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Text, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	ImageMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Image, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	VoiceMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Voice, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	VideoMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Video, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	LocationMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Location, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)

	// 事件处理函数
	//  rawXMLMsg 是解密后的明文 xml 消息体
	//  timestamp 是请求中的时间戳
	//  nonce     是请求中的随机数
	//  random    是请求中的消息体加密的 random
	SubscribeEventHandler(w http.ResponseWriter, r *http.Request, event *request.SubscribeEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	UnsubscribeEventHandler(w http.ResponseWriter, r *http.Request, event *request.UnsubscribeEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	LocationEventHandler(w http.ResponseWriter, r *http.Request, event *request.LocationEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	MenuClickEventHandler(w http.ResponseWriter, r *http.Request, event *request.MenuClickEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
	MenuViewEventHandler(w http.ResponseWriter, r *http.Request, event *request.MenuViewEvent, rawXMLMsg []byte, timestamp int64, nonce string, random []byte)
}
