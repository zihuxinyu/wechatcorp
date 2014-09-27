刚写好, 还没有测试, 哪位仁兄测试下
```golang
package main

import (
	"github.com/chanxuehong/wechatcorp/message/passive/request"
	"github.com/chanxuehong/wechatcorp/message/passive/response"
	"github.com/chanxuehong/wechatcorp/server"
	"net/http"
)

// 自己实现一个 server.AgentMsgHandler
type CustomAgentMsgHandler struct {
	server.DefaultAgentMsgHandler // 提供了默认实现
}

// 自定义文本消息处理函数, 覆盖默认的实现
func (handler *CustomAgentMsgHandler) TextMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Text, rawXMLMsg []byte, timestamp, nonce int64, random []byte) {
	// 示例代码, 把用户发送过来的文本原样的回复过去

	w.Header().Set("Content-Type", "application/xml; charset=utf-8") // 可选
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.Content)
	if err := handler.WriteText(w, resp, timestamp, nonce, random); err != nil {
		// TODO: 增加错误处理代码
	}
}

func init() {
	// TODO: 获取必要数据的代码

	var AgentMsgHandler CustomAgentMsgHandler
	AgentMsgHandler.Init("CorpId", "AgentId", "Token", []byte("AESKey"))

	var AgentHttpHandler server.AgentHttpHandler
	AgentHttpHandler.AgentMsgHandler = &AgentMsgHandler

	// 注册这个 handler 到回调 URL 上
	// 比如你在公众平台后台注册的回调地址是 http://abc.xxx.com/weixin，那么可以这样注册
	http.Handle("/weixin", AgentHttpHandler)
}

func main() {
	http.ListenAndServe(":80", nil)
}
```