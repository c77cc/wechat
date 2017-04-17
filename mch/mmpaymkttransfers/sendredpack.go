// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package mmpaymkttransfers

import (
	"github.com/c77cc/wechat/mch"
)

// 红包发放.
//  NOTE: 请求需要双向证书
func SendRedPack(proxy *mch.Proxy, req map[string]string) (resp map[string]string, err error) {
	return proxy.PostXML("https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack", req)
}
