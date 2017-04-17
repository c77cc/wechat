// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package component

import (
	"net/url"
)

// 请求用户授权时跳转的地址.
func AuthCodeURL(appId, preAuthCode, redirectURI, state string) string {
	return "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=" + url.QueryEscape(appId) +
		"&pre_auth_code=" + url.QueryEscape(preAuthCode) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&state=" + url.QueryEscape(state)
}
