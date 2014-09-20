// @description wechat-corp 是腾讯微信公众平台企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package request

type SubscribeEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonHead

	Event string `xml:"Event" json:"Event"` // 事件类型，subscribe(订阅)
}

type UnsubscribeEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonHead

	Event string `xml:"Event" json:"Event"` // 事件类型，unsubscribe(取消订阅)
}

type LocationEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonHead

	Event     string  `xml:"Event"     json:"Event"`     // 事件类型，LOCATION
	Latitude  float64 `xml:"Latitude"  json:"Latitude"`  // 地理位置纬度
	Longitude float64 `xml:"Longitude" json:"Longitude"` // 地理位置经度
	Precision float64 `xml:"Precision" json:"Precision"` // 地理位置精度
}

type MenuClickEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonHead

	Event    string `xml:"Event"    json:"Event"`    // 事件类型，CLICK
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值，与自定义菜单接口中KEY值对应
}

type MenuViewEvent struct {
	XMLName struct{} `xml:"xml" json:"-"`
	CommonHead

	Event    string `xml:"Event"    json:"Event"`    // 事件类型，VIEW
	EventKey string `xml:"EventKey" json:"EventKey"` // 事件KEY值，设置的跳转URL
}
