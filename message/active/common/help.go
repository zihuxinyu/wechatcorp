// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package common

import (
	"strconv"
	"strings"
)

// 用 '|' 连接 a 的各个元素
func JoinString(a []string) string {
	return strings.Join(a, "|")
}

// 用 '|' 连接 a 的各个元素的十进制字符串
func JoinInt64(a []int64) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return strconv.FormatInt(a[0], 10)
	}

	b := make([]string, len(a))
	for i, n := range a {
		b[i] = strconv.FormatInt(n, 10)
	}

	return strings.Join(b, "|")
}
