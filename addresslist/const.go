// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package addresslist

const (
	USERINFO_GENDER_MALE   = 0 // 男性
	USERINFO_GENDER_FEMALE = 1 // 女性

	USERINFO_ENABLE_TRUE  = 1 // 启用成员
	USERINFO_ENABLE_FALSE = 0 // 禁用成员

	USERINFO_STATUS_SUBSCRIBED   = 1 // 已关注
	USERINFO_STATUS_BLOCKED      = 2 // 已冻结
	USERINFO_STATUS_NOSUBSCRIBED = 3 // 未关注
)
