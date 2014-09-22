// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 获取缓存中的 access token, err == nil 时 token 才有效!
func (c *Client) Token() (token string, err error) {
	if token, err = c.tokenCache.Token(); err != ErrCacheMiss {
		return
	}
	return c.TokenRefresh()
}

// 从微信服务器获取有效的 access token 并更新 TokenCache, err == nil 时 token 才有效!
//  NOTE: 一般情况下无需调用该函数, 请使用 Token() 获取 access token.
func (c *Client) TokenRefresh() (token string, err error) {
	if token, err = c.getToken(); err != nil {
		return
	}
	if err = c.tokenCache.SetToken(token); err != nil {
		return
	}
	return
}

// 从微信服务器获取新的 access_token
func (c *Client) getToken() (token string, err error) {
	_url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" +
		c.corpId + "&corpsecret=" + c.corpSecret

	httpResp, err := c.httpClient.Get(_url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	var result struct {
		Error
		Token string `json:"access_token"`
	}
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	if result.ErrCode != errCodeOK {
		err = &result.Error
		return
	}

	token = result.Token
	return
}
