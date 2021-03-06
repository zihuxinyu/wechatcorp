// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package wechatcorp

import (
	"errors"
	"strconv"
	"strings"
)

// 获取微信客户端的版本.
//  @userAgent: 微信内置浏览器的 user-agent;
//  @x, y, z:   如果微信版本为 5.3.1 则有 x==5, y==3, z==1
//  @err:       错误信息
func WXVersion(userAgent string) (x, y, z int, err error) {
	// Mozilla/5.0(iphone;CPU iphone OS 5_1_1 like Mac OS X) AppleWebKit/534.46(KHTML,like Geocko)Mobile/9B206 MicroMessenger/5.0
	lastSlashIndex := strings.LastIndex(userAgent, "/")
	versionIndex := lastSlashIndex + 1

	if lastSlashIndex == -1 || versionIndex == len(userAgent) {
		err = errors.New("不是有效的微信浏览器 user-agent")
		return
	}

	strArr := strings.Split(userAgent[versionIndex:], ".")
	verArr := make([]int, len(strArr))

	// strArr 的每个元素应该都是整数
	for i, str := range strArr {
		var ver uint64
		ver, err = strconv.ParseUint(str, 10, 16)
		if err != nil {
			err = errors.New("不是有效的微信浏览器 user-agent")
			return
		}
		verArr[i] = int(ver)
	}

	// verArr 至少有一个元素
	switch len(verArr) {
	case 3:
		x = verArr[0]
		y = verArr[1]
		z = verArr[2]

	case 2:
		x = verArr[0]
		y = verArr[1]

	case 1:
		x = verArr[0]

	default:
		x = verArr[0]
		y = verArr[1]
		z = verArr[2]
	}
	return
}
