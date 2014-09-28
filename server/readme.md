```golang
package main

import (
	"github.com/chanxuehong/wechatcorp/message/passive/request"
	"github.com/chanxuehong/wechatcorp/message/passive/response"
	"github.com/chanxuehong/wechatcorp/server"
	"log"
	"net/http"
	"time"
)

// 自己实现一个 server.AgentMsgHandler
type CustomAgentMsgHandler struct {
	server.DefaultAgentMsgHandler // 可选! 提供了默认实现
}

// 自定义文本消息处理函数, 覆盖默认的实现
func (handler *CustomAgentMsgHandler) TextMsgHandler(w http.ResponseWriter, r *http.Request, msg *request.Text, rawXMLMsg []byte, timestamp int64, nonce string, random []byte) {
	// 示例代码
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")                           // 可选
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.Content, time.Now().Unix()) // 时间戳也可以用传入的参数 timestamp

	// timestamp, nonce, random 参数可以直接用传入的参数, 也可以自己生成!!!
	if err := handler.WriteText(w, resp, timestamp, nonce, random); err != nil {
		// TODO: 增加错误处理代码
	}
}

func CustomInvalidRequestHandlerFunc(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
}

func init() {
	// TODO: 获取必要数据的代码

	var AgentMsgHandler CustomAgentMsgHandler
	AgentMsgHandler.DefaultAgentMsgHandler.Init("CorpId", "AgentId", "Token", []byte("AESKey"))

	var HttpHandler server.NCSHttpHandler
	HttpHandler.SetInvalidRequestHandler(server.InvalidRequestHandlerFunc(CustomInvalidRequestHandlerFunc))
	HttpHandler.SetAgentMsgHandler("CorpId", "AgentId", &AgentMsgHandler)

	// 注册这个 handler 到回调 URL 上
	// 比如你在公众平台后台注册的回调地址是 http://abc.xxx.com/weixin，那么可以这样注册
	http.Handle("/weixin", &HttpHandler)
}

func main() {
	http.ListenAndServe(":80", nil)
}
```