// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package addresslist

type Department struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ParentId int64  `json:"parentid"`
}
