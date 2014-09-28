// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

// 把整数 n 格式化成 4 字节的网络字节序
func encodeNetworkBytesOrder(n int, orderBytes []byte) {
	if len(orderBytes) != 4 {
		panic("the length of orderBytes must be equal to 4")
	}
	orderBytes[0] = byte(n >> 24)
	orderBytes[1] = byte(n >> 16)
	orderBytes[2] = byte(n >> 8)
	orderBytes[3] = byte(n)
}

// 从 4 字节的网络字节序里解析出整数
func decodeNetworkBytesOrder(orderBytes []byte) (n int) {
	if len(orderBytes) != 4 {
		panic("the length of orderBytes must be equal to 4")
	}
	n = int(orderBytes[0])<<24 |
		int(orderBytes[1])<<16 |
		int(orderBytes[2])<<8 |
		int(orderBytes[3])
	return
}

func (handler *DefaultAgentMsgHandler) EncryptMsg(random, rawXMLMsg []byte) (EncryptMsg []byte) {
	buf := make([]byte, 20+len(rawXMLMsg)+len(handler.CorpId)+aes.BlockSize)
	plain := buf[:20]
	pad := buf[len(buf)-aes.BlockSize:]

	// 拼接
	copy(plain, random) // 使用参数 random, 不自己生成
	encodeNetworkBytesOrder(len(rawXMLMsg), plain[16:20])
	plain = append(plain, rawXMLMsg...)
	plain = append(plain, handler.CorpId...)

	// 补位
	amountToPad := aes.BlockSize - len(plain)%aes.BlockSize
	pad = pad[:amountToPad]
	for i := 0; i < amountToPad; i++ {
		pad[i] = byte(amountToPad)
	}
	plain = append(plain, pad...)

	// 加密
	block, err := aes.NewCipher(handler.AESKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, handler.CipherIV)
	mode.CryptBlocks(plain, plain)

	EncryptMsg = plain
	return
}

func (handler *DefaultAgentMsgHandler) DecryptMsg(EncryptMsg []byte) (random, rawXMLMsg []byte, err error) {
	// 解密
	if len(EncryptMsg) < aes.BlockSize {
		err = errors.New("EncryptMsg too short")
		return
	}
	if len(EncryptMsg)%aes.BlockSize != 0 {
		err = errors.New("EncryptMsg is not a multiple of the block size")
		return
	}

	block, err := aes.NewCipher(handler.AESKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, handler.CipherIV)
	mode.CryptBlocks(EncryptMsg, EncryptMsg)

	plain := EncryptMsg

	// 去除补位
	amountToPad := int(plain[len(plain)-1])
	if amountToPad < 1 || amountToPad > aes.BlockSize {
		err = errors.New("the amount to pad is invalid")
		return
	}
	plain = plain[:len(plain)-amountToPad]

	// 反拼装
	if len(plain) <= 20 {
		err = errors.New("plain too short")
		return
	}
	msgLen := decodeNetworkBytesOrder(plain[16:20])
	if msgLen < 0 {
		err = fmt.Errorf("invalid msg length: %d", msgLen)
		return
	}
	msgEnd := 20 + msgLen
	if msgEnd >= len(plain) {
		err = fmt.Errorf("msg length is too large: %d", msgLen)
		return
	}

	CorpId := string(plain[msgEnd:])
	if CorpId != handler.CorpId {
		err = errors.New("CorpId mismatch")
		return
	}

	random = plain[:16]
	rawXMLMsg = plain[20:msgEnd]
	return
}
