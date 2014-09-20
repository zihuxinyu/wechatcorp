// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package menu

const (
	MenuButtonCountLimit    = 3 // 一级菜单最多包含 3 个按钮
	SubMenuButtonCountLimit = 5 // 二级菜单最多包含 5 个按钮
)

const (
	MenuButtonLenLimit    = 16 // 菜单标题不超过16个字节
	SubMenuButtonLenLimit = 40 // 子菜单标题不超过40个字节
)
const (
	ButtonKeyLenLimit = 128 // 菜单KEY值不能超过128字节
	ButtonURLLenLimit = 256 // 网页链接不能超过256字节
)

const (
	BUTTON_TYPE_CLICK = "click"
	BUTTON_TYPE_VIEW  = "view"
)
