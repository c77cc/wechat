// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

// +build wechatdebug

package mch

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/c77cc/util"
)

type Proxy struct {
	apiKey     string
	httpClient *http.Client
}

// 创建一个新的 Proxy.
//  如果 httpClient == nil 则默认用 http.DefaultClient.
func NewProxy(apiKey string, httpClient *http.Client) *Proxy {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Proxy{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// 微信支付通用请求方法.
//  注意: err == nil 表示协议状态都为 SUCCESS.
func (proxy *Proxy) PostXML(url string, req map[string]string) (resp map[string]string, err error) {
	bodyBuf := textBufferPool.Get().(*bytes.Buffer)
	bodyBuf.Reset()
	defer textBufferPool.Put(bodyBuf)

	if err = util.FormatMapToXML(bodyBuf, req); err != nil {
		return
	}

	LogInfoln("[WECHAT_DEBUG] request url:", url)
	LogInfoln("[WECHAT_DEBUG] request xml:", bodyBuf.String())

	httpResp, err := proxy.httpClient.Post(url, "text/xml; charset=utf-8", bodyBuf)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	LogInfoln("[WECHAT_DEBUG] response xml:", string(respBody))

	if resp, err = util.ParseXMLToMap(bytes.NewReader(respBody)); err != nil {
		return
	}

	// 判断协议状态
	ReturnCode, ok := resp["return_code"]
	if !ok {
		err = errors.New("no return_code parameter")
		return
	}
	if ReturnCode != ReturnCodeSuccess {
		err = &Error{
			ReturnCode: ReturnCode,
			ReturnMsg:  resp["return_msg"],
		}
		return
	}

	// 认证签名
	signature1, ok := resp["sign"]
	if !ok {
		err = errors.New("no sign parameter")
		return
	}
	signature2 := Sign(resp, proxy.apiKey, nil)
	if signature1 != signature2 {
		err = fmt.Errorf("check signature failed, \r\ninput: %q, \r\nlocal: %q", signature1, signature2)
		return
	}
	return
}
