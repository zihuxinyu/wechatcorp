// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package request

type CommonHead struct {
	ToUserName   string `xml:"ToUserName"   json:"ToUserName"`   // 企业号CorpID
	FromUserName string `xml:"FromUserName" json:"FromUserName"` // 员工UserID
	CreateTime   int64  `xml:"CreateTime"   json:"CreateTime"`   // 消息创建时间（整型）, unixtime
	MsgType      string `xml:"MsgType"      json:"MsgType"`      // 消息类型
	AgentID      string `xml:"AgentID"      json:"AgentID"`      // 企业应用的id
}
