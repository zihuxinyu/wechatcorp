// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package request

// 微信服务器请求 http body
//
//  <xml>
//      <ToUserName><![CDATA[toUser]]></ToUserName>
//      <AgentID><![CDATA[toAgentID]]></AgentID>
//      <Encrypt><![CDATA[msg_encrypt]]></Encrypt>
//  </xml>
type RequestHttpBody struct {
	XMLName    struct{} `xml:"xml" json:"-"`
	CorpId     string   `xml:"ToUserName"`
	AgentId    int64    `xml:"AgentID"`
	EncryptMsg string   `xml:"Encrypt"` // EncryptMsg 为经过加密的密文
}
