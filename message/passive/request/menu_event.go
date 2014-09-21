// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package request

// 点击菜单拉取消息时的事件推送
type MenuClickEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonHead

	Event    string `xml:"Event"    json:"Event"`    // 事件类型，CLICK
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值，与自定义菜单接口中KEY值对应
}

func (req *Request) MenuClickEvent() (event *MenuClickEvent) {
	event = &MenuClickEvent{
		CommonHead: req.CommonHead,
		Event:      req.Event,
		EventKey:   req.EventKey,
	}
	return
}

// 点击菜单跳转链接时的事件推送
type MenuViewEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonHead

	Event    string `xml:"Event"    json:"Event"`    // 事件类型，VIEW
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值，设置的跳转URL
}

func (req *Request) MenuViewEvent() (event *MenuViewEvent) {
	event = &MenuViewEvent{
		CommonHead: req.CommonHead,
		Event:      req.Event,
		EventKey:   req.EventKey,
	}
	return
}
