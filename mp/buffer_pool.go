// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package mp

import (
	"bytes"
	"sync"
)

var textBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 16<<10)) // 16KB
	},
}

var mediaBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 10<<20)) // 10MB
	},
}
