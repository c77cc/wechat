// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package mch

import (
	"net/http"
)

// 微信服务器推送过来的消息(事件)处理接口
type MessageHandler interface {
	ServeMessage(w http.ResponseWriter, r *Request)
}

type MessageHandlerFunc func(http.ResponseWriter, *Request)

func (fn MessageHandlerFunc) ServeMessage(w http.ResponseWriter, r *Request) {
	fn(w, r)
}

// 消息(事件)请求信息
type Request struct {
	HttpRequest *http.Request // 可以为 nil, 因为某些 http 框架没有提供此参数

	RawMsgXML []byte            // 消息的 XML 文本
	Msg       map[string]string // 解析后的消息
}
