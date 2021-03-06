// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package component

import (
	"errors"
	"sync"
)

type Server interface {
	AppId() string // 获取第三方平台AppId
	Token() string // 获取第三方平台的Token

	CurrentAESKey() [32]byte // 获取当前有效的 AES 加密 Key
	LastAESKey() [32]byte    // 获取最后一个有效的 AES 加密 Key

	MessageHandler() MessageHandler // 获取 MessageHandler
}

var _ Server = (*DefaultServer)(nil)

type DefaultServer struct {
	appId string
	token string

	rwmutex           sync.RWMutex
	currentAESKey     [32]byte // 当前的 AES Key
	lastAESKey        [32]byte // 最后一个 AES Key
	isLastAESKeyValid bool     // lastAESKey 是否有效, 如果 lastAESKey 是 zero 则无效

	messageHandler MessageHandler
}

// NewDefaultServer 创建一个新的 DefaultServer.
func NewDefaultServer(appId, token string, AESKey []byte, handler MessageHandler) (srv *DefaultServer) {
	if len(AESKey) != 32 {
		panic("the length of AESKey must equal to 32")
	}
	if handler == nil {
		panic("nil MessageHandler")
	}

	srv = &DefaultServer{
		appId:          appId,
		token:          token,
		messageHandler: handler,
	}
	copy(srv.currentAESKey[:], AESKey)
	return
}

func (srv *DefaultServer) AppId() string {
	return srv.appId
}
func (srv *DefaultServer) Token() string {
	return srv.token
}
func (srv *DefaultServer) MessageHandler() MessageHandler {
	return srv.messageHandler
}
func (srv *DefaultServer) CurrentAESKey() (key [32]byte) {
	srv.rwmutex.RLock()
	key = srv.currentAESKey
	srv.rwmutex.RUnlock()
	return
}
func (srv *DefaultServer) LastAESKey() (key [32]byte) {
	srv.rwmutex.RLock()
	if srv.isLastAESKeyValid {
		key = srv.lastAESKey
	} else {
		key = srv.currentAESKey
	}
	srv.rwmutex.RUnlock()
	return
}
func (srv *DefaultServer) UpdateAESKey(AESKey []byte) (err error) {
	if len(AESKey) != 32 {
		return errors.New("the length of AESKey must equal to 32")
	}

	srv.rwmutex.Lock()
	srv.isLastAESKeyValid = true
	srv.lastAESKey = srv.currentAESKey
	copy(srv.currentAESKey[:], AESKey)
	srv.rwmutex.Unlock()
	return
}
