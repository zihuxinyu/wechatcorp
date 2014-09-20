// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package media

// 上传(创建)媒体成功时的回复报文
type MediaInfo struct {
	MediaType string `json:"type"` // 图片（image）、语音（voice）、视频（video）、普通文件（file）
	MediaId   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
}
