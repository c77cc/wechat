// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package thirdparty

import (
	"bytes"
	"crypto/subtle"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/c77cc/wechat/corp"
	"github.com/c77cc/wechat/util"
)

var zeroAESKey [32]byte

// 微信服务器请求 http body
type RequestHttpBody struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	SuiteId      string   `xml:"ToUserName"`
	EncryptedMsg string   `xml:"Encrypt"`
}

// ServeHTTP 处理 http 消息请求
//  NOTE: 调用者保证所有参数有效
func ServeHTTP(w http.ResponseWriter, r *http.Request, queryValues url.Values,
	suiteServer SuiteServer, invalidRequestHandler corp.InvalidRequestHandler) {

	switch r.Method {
	case "POST": // 消息处理
		msgSignature1, timestampStr, nonce, err := parsePostURLQuery(queryValues)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// 首先判断签名长度是否合法
		if len(msgSignature1) != 40 {
			err = fmt.Errorf("the length of msg_signature mismatch, have: %d, want: 40", len(msgSignature1))
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			err = errors.New("can not parse timestamp to int64: " + timestampStr)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// 解析 RequestHttpBody
		var requestHttpBody RequestHttpBody
		if err := xml.NewDecoder(r.Body).Decode(&requestHttpBody); err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		haveSuiteId := requestHttpBody.SuiteId
		wantSuiteId := suiteServer.SuiteId()
		if len(haveSuiteId) != len(wantSuiteId) {
			err = fmt.Errorf("the RequestHttpBody's ToUserName mismatch, have: %s, want: %s", haveSuiteId, wantSuiteId)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}
		if subtle.ConstantTimeCompare([]byte(haveSuiteId), []byte(wantSuiteId)) != 1 {
			err = fmt.Errorf("the RequestHttpBody's ToUserName mismatch, have: %s, want: %s", haveSuiteId, wantSuiteId)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		suiteToken := suiteServer.SuiteToken()

		// 验证签名
		msgSignature2 := util.MsgSign(suiteToken, timestampStr, nonce, requestHttpBody.EncryptedMsg)
		if subtle.ConstantTimeCompare([]byte(msgSignature1), []byte(msgSignature2)) != 1 {
			err = fmt.Errorf("check signature failed, input: %s, local: %s", msgSignature1, msgSignature2)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// 解密
		EncryptedMsgBytes, err := base64.StdEncoding.DecodeString(requestHttpBody.EncryptedMsg)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		AESKey := suiteServer.CurrentAESKey()
		Random, RawMsgXML, err := util.AESDecryptMsg(EncryptedMsgBytes, wantSuiteId, AESKey)
		if err != nil {
			// 尝试用上一次的 AESKey 来解密
			LastAESKey := suiteServer.LastAESKey()
			if bytes.Equal(AESKey[:], LastAESKey[:]) || bytes.Equal(zeroAESKey[:], LastAESKey[:]) {
				invalidRequestHandler.ServeInvalidRequest(w, r, err)
				return
			}

			AESKey = LastAESKey // NOTE
			Random, RawMsgXML, err = util.AESDecryptMsg(EncryptedMsgBytes, wantSuiteId, AESKey)
			if err != nil {
				invalidRequestHandler.ServeInvalidRequest(w, r, err)
				return
			}
		}

		// 解密成功, 解析 MixedMessage
		var MixedMsg MixedSuiteMessage
		if err = xml.Unmarshal(RawMsgXML, &MixedMsg); err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// 安全考虑再次验证
		if haveSuiteId != MixedMsg.SuiteId {
			err = fmt.Errorf("the RequestHttpBody's ToUserName(==%s) mismatch the MixedMessage's SuiteId(==%s)", haveSuiteId, MixedMsg.SuiteId)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// 成功, 交给 SuiteMessageHandler
		r := &Request{
			HttpRequest: r,

			QueryValues:  queryValues,
			MsgSignature: msgSignature1,
			Timestamp:    timestamp,
			Nonce:        nonce,

			RawMsgXML: RawMsgXML,
			MixedMsg:  &MixedMsg,

			AESKey: AESKey,
			Random: Random,

			SuiteId:    haveSuiteId,
			SuiteToken: suiteToken,
		}
		suiteServer.SuiteMessageHandler().ServeMessage(w, r)

	case "GET": // 首次验证
		msgSignature1, timestamp, nonce, encryptedMsg, err := parseGetURLQuery(queryValues)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// 验证签名
		if len(msgSignature1) != 40 {
			err = fmt.Errorf("the length of msg_signature mismatch, have: %d, want: 40", len(msgSignature1))
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		msgSignature2 := util.MsgSign(suiteServer.SuiteToken(), timestamp, nonce, encryptedMsg)
		if subtle.ConstantTimeCompare([]byte(msgSignature1), []byte(msgSignature2)) != 1 {
			err = fmt.Errorf("check signature failed, input: %s, local: %s", msgSignature1, msgSignature2)
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		// 解密
		EncryptedMsgBytes, err := base64.StdEncoding.DecodeString(encryptedMsg)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		SuiteId := suiteServer.SuiteId()
		AESKey := suiteServer.CurrentAESKey()
		_, echostr, err := util.AESDecryptMsg(EncryptedMsgBytes, SuiteId, AESKey)
		if err != nil {
			invalidRequestHandler.ServeInvalidRequest(w, r, err)
			return
		}

		w.Write(echostr)
	}
}
