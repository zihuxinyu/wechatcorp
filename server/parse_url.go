// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package server

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
)

func parsePostURLQuery(URL *url.URL) (signature, timestamp, nonce string, err error) {
	if URL == nil {
		err = errors.New("URL == nil")
		return
	}

	urlValues, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return
	}

	signature = urlValues.Get("msg_signature")
	if signature == "" {
		err = errors.New("msg_signature is empty")
		return
	}

	const signatureLen = sha1.Size * 2
	if len(signature) != signatureLen {
		err = fmt.Errorf("the length of msg_signature mismatch, have: %d, want: %d", len(signature), signatureLen)
		return
	}

	timestamp = urlValues.Get("timestamp")
	if timestamp == "" {
		err = errors.New("timestamp is empty")
		return
	}

	nonce = urlValues.Get("nonce")
	if nonce == "" {
		err = errors.New("nonce is empty")
		return
	}

	return
}

func parsePostURLQueryEx(URL *url.URL) (agentkey, signature, timestamp, nonce string, err error) {
	if URL == nil {
		err = errors.New("URL == nil")
		return
	}

	urlValues, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return
	}

	agentkey = urlValues.Get(URLQueryAgentKeyName)
	if agentkey == "" {
		err = errors.New(URLQueryAgentKeyName + " is empty")
		return
	}

	signature = urlValues.Get("msg_signature")
	if signature == "" {
		err = errors.New("msg_signature is empty")
		return
	}

	const signatureLen = sha1.Size * 2
	if len(signature) != signatureLen {
		err = fmt.Errorf("the length of msg_signature mismatch, have: %d, want: %d", len(signature), signatureLen)
		return
	}

	timestamp = urlValues.Get("timestamp")
	if timestamp == "" {
		err = errors.New("timestamp is empty")
		return
	}

	nonce = urlValues.Get("nonce")
	if nonce == "" {
		err = errors.New("nonce is empty")
		return
	}

	return
}

func parseGetURLQuery(URL *url.URL) (signature, timestamp, nonce, echostr string, err error) {
	if URL == nil {
		err = errors.New("URL == nil")
		return
	}

	urlValues, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return
	}

	signature = urlValues.Get("msg_signature")
	if signature == "" {
		err = errors.New("msg_signature is empty")
		return
	}

	const signatureLen = sha1.Size * 2
	if len(signature) != signatureLen {
		err = fmt.Errorf("the length of msg_signature mismatch, have: %d, want: %d", len(signature), signatureLen)
		return
	}

	timestamp = urlValues.Get("timestamp")
	if timestamp == "" {
		err = errors.New("timestamp is empty")
		return
	}

	nonce = urlValues.Get("nonce")
	if nonce == "" {
		err = errors.New("nonce is empty")
		return
	}

	echostr = urlValues.Get("echostr")
	if echostr == "" {
		err = errors.New("echostr is empty")
		return
	}

	return
}

func parseGetURLQueryEx(URL *url.URL) (agentkey, signature, timestamp, nonce, echostr string, err error) {
	if URL == nil {
		err = errors.New("URL == nil")
		return
	}

	urlValues, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return
	}

	agentkey = urlValues.Get(URLQueryAgentKeyName)
	if agentkey == "" {
		err = errors.New(URLQueryAgentKeyName + " is empty")
		return
	}

	signature = urlValues.Get("msg_signature")
	if signature == "" {
		err = errors.New("msg_signature is empty")
		return
	}

	const signatureLen = sha1.Size * 2
	if len(signature) != signatureLen {
		err = fmt.Errorf("the length of msg_signature mismatch, have: %d, want: %d", len(signature), signatureLen)
		return
	}

	timestamp = urlValues.Get("timestamp")
	if timestamp == "" {
		err = errors.New("timestamp is empty")
		return
	}

	nonce = urlValues.Get("nonce")
	if nonce == "" {
		err = errors.New("nonce is empty")
		return
	}

	echostr = urlValues.Get("echostr")
	if echostr == "" {
		err = errors.New("echostr is empty")
		return
	}

	return
}
