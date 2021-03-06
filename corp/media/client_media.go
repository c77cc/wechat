// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"

	"github.com/c77cc/wechat/corp"
)

// 下载多媒体到文件.
func (clt Client) DownloadMedia(mediaId, filepath string) (err error) {
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer func() {
		file.Close()
		if err != nil {
			os.Remove(filepath)
		}
	}()

	return clt.downloadMediaToWriter(mediaId, file)
}

// 下载多媒体到 io.Writer.
func (clt Client) DownloadMediaToWriter(mediaId string, writer io.Writer) error {
	if writer == nil {
		return errors.New("nil writer")
	}
	return clt.downloadMediaToWriter(mediaId, writer)
}

// 下载多媒体到 io.Writer.
func (clt Client) downloadMediaToWriter(mediaId string, writer io.Writer) (err error) {
	token, err := clt.Token()
	if err != nil {
		return
	}

	hasRetried := false
RETRY:
	finalURL := "https://qyapi.weixin.qq.com/cgi-bin/media/get?media_id=" + url.QueryEscape(mediaId) +
		"&access_token=" + url.QueryEscape(token)

	httpResp, err := clt.HttpClient.Get(finalURL)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	ContentType, _, _ := mime.ParseMediaType(httpResp.Header.Get("Content-Type"))
	if ContentType != "text/plain" && ContentType != "application/json" {
		// 返回的是媒体流
		_, err = io.Copy(writer, httpResp.Body)
		return
	}

	// 返回的是错误信息
	var result corp.Error
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return
	}

	switch result.ErrCode {
	case corp.ErrCodeOK:
		return // 基本不会出现
	case corp.ErrCodeAccessTokenExpired: // 失效(过期)重试一次
		corp.LogInfoln("[WECHAT_RETRY] err_code:", result.ErrCode, ", err_msg:", result.ErrMsg)
		corp.LogInfoln("[WECHAT_RETRY] current token:", token)

		if !hasRetried {
			hasRetried = true

			if token, err = clt.TokenRefresh(); err != nil {
				return
			}
			corp.LogInfoln("[WECHAT_RETRY] new token:", token)

			result = corp.Error{}
			goto RETRY
		}
		corp.LogInfoln("[WECHAT_RETRY] fallthrough, current token:", token)
		fallthrough
	default:
		err = &result
		return
	}
}
