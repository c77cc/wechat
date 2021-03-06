// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package thirdparty

import (
	"github.com/c77cc/wechat/corp"
)

type PreAuthCode struct {
	Value     string `json:"pre_auth_code"`
	ExpiresIn int64  `json:"expires_in"`
}

// 获取预授权码.
//  AppIdList: 应用id，本参数选填，表示用户能对本套件内的哪些应用授权，不填时默认用户有全部授权权限
func (clt *SuiteClient) GetPreAuthCode(AppIdList []int64) (code *PreAuthCode, err error) {
	request := struct {
		SuiteId   string  `json:"suite_id"`
		AppIdList []int64 `json:"appid,omitempty"`
	}{
		SuiteId:   clt.SuiteId,
		AppIdList: AppIdList,
	}

	var result struct {
		corp.Error
		PreAuthCode
	}

	incompleteURL := "https://qyapi.weixin.qq.com/cgi-bin/service/get_pre_auth_code?suite_access_token="
	if err = clt.PostJSON(incompleteURL, &request, &result); err != nil {
		return
	}

	if result.ErrCode != corp.ErrCodeOK {
		err = &result.Error
		return
	}
	code = &result.PreAuthCode
	return
}
