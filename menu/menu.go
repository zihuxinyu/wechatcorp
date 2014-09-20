// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package menu

// 菜单
//
//  {
//      "button": [
//          {
//              "type": "click",
//              "name": "今日歌曲",
//              "key": "V1001_TODAY_MUSIC"
//          },
//          {
//              "name": "菜单",
//              "sub_button": [
//                  {
//                      "type": "view",
//                      "name": "搜索",
//                      "url": "http://www.soso.com/"
//                  },
//                  {
//                      "type": "click",
//                      "name": "赞一下我们",
//                      "key": "V1001_GOOD"
//                  }
//              ]
//          }
//      ]
//  }
//
type Menu struct {
	Buttons []Button `json:"button,omitempty"` // 按钮个数不能超过 MenuButtonCountLimit
}

// 菜单的按钮
type Button struct {
	Name       string   `json:"name"`                 // 菜单标题，不超过16个字节，子菜单不超过40个字节
	Type       string   `json:"type,omitempty"`       // 菜单的响应动作类型，目前有click、view两种类型
	Key        string   `json:"key,omitempty"`        // click类型必须; 菜单KEY值，用于消息接口推送，不超过128字节
	URL        string   `json:"url,omitempty"`        // view类型必须; 网页链接，员工点击菜单可打开链接，不超过256字节
	SubButtons []Button `json:"sub_button,omitempty"` // 二级菜单, 按钮个数不能超过 SubMenuButtonCountLimit
}

// 初始化 btn 指向的 Button 为 click 类型按钮
func (btn *Button) InitToClickButton(name, key string) {
	btn.Name = name
	btn.Type = BUTTON_TYPE_CLICK
	btn.Key = key

	// 容错性考虑, 清除其他字段
	btn.URL = ""
	btn.SubButtons = nil
}

// 初始化 btn 指向的 Button 为 view 类型按钮
func (btn *Button) InitToViewButton(name, url string) {
	btn.Name = name
	btn.Type = BUTTON_TYPE_VIEW
	btn.URL = url

	// 容错性考虑, 清除其他字段
	btn.Key = ""
	btn.SubButtons = nil
}

// 初始化 btn 指向的 Button 为 子菜单 类型按钮
func (btn *Button) InitToSubMenuButton(name string, subButtons []Button) {
	btn.Name = name
	btn.SubButtons = subButtons

	// 容错性考虑, 清除其他字段
	btn.Type = ""
	btn.Key = ""
	btn.URL = ""
}
