// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package response

// 回复微信请求的 http body
//
//  <xml>
//      <Encrypt><![CDATA[msg_encrypt]]></Encrypt>
//      <MsgSignature><![CDATA[msg_signature]]></MsgSignature>
//      <TimeStamp>timestamp</TimeStamp>
//      <Nonce><![CDATA[nonce]]></Nonce>
//  </xml>
type ResponseHttpBody struct {
	XMLName    struct{} `xml:"xml" json:"-"`
	EncryptMsg string   `xml:"Encrypt"` // EncryptMsg 为经过加密的密文
	Signature  string   `xml:"MsgSignature"`
	TimeStamp  int64    `xml:"TimeStamp"`
	Nonce      string   `xml:"Nonce"`
}
