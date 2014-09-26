// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"crypto/sha1"
	"encoding/hex"
	"sort"
)

func (handler *DefaultAgentMsgHandler) Signature(timestamp, nonce, EncryptMsg string) (signature string) {
	strArray := sort.StringSlice{timestamp, nonce, EncryptMsg, handler.Token}
	strArray.Sort()

	n := len(timestamp) + len(nonce) + len(EncryptMsg) + len(handler.Token)
	buf := make([]byte, 0, n)

	buf = append(buf, strArray[0]...)
	buf = append(buf, strArray[1]...)
	buf = append(buf, strArray[2]...)
	buf = append(buf, strArray[3]...)

	hashSumArray := sha1.Sum(buf)
	return hex.EncodeToString(hashSumArray[:])
}
