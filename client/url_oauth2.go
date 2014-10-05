// @description wechatcorp 是腾讯微信公众平台 企业号 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechatcorp for the canonical source repository
// @license     https://github.com/chanxuehong/wechatcorp/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package client

// https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=ACCESS_TOKEN
// &code=CODE&agentid=AGENTID
func _OAuth2GetUserInfoURL(accesstoken string, code string, agentid string) string {
	return "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=" +
		accesstoken + "&code=" + code + "&agentid=" + agentid
}

// https://qyapi.weixin.qq.com/cgi-bin/user/authsucc?access_token=ACCESS_TOKEN&userid=USERID
func _OAuth2UserAuthSuccessfullyURL(accesstoken string, userid string) string {
	return "https://qyapi.weixin.qq.com/cgi-bin/user/authsucc?access_token=" +
		accesstoken + "&userid=" + userid
}
