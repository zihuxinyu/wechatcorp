// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package request

const (
	// 微信服务器推送过来的消息类型
	MSG_TYPE_TEXT     = "text"
	MSG_TYPE_IMAGE    = "image"
	MSG_TYPE_VOICE    = "voice"
	MSG_TYPE_VIDEO    = "video"
	MSG_TYPE_LOCATION = "location"
	MSG_TYPE_EVENT    = "event"
)

const (
	// 微信服务器推送过来的事件类型
	EVENT_TYPE_SUBSCRIBE   = "subscribe"   // 关注
	EVENT_TYPE_UNSUBSCRIBE = "unsubscribe" // 取消关注
	EVENT_TYPE_CLICK       = "CLICK"       // 点击菜单拉取消息时的事件推送
	EVENT_TYPE_VIEW        = "VIEW"        // 点击菜单跳转链接时的事件推送
	EVENT_TYPE_LOCATION    = "LOCATION"    // 上报地理位置事件
)
