// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/c77cc/wechat for the canonical source repository
// @license     https://github.com/c77cc/wechat/blob/master/LICENSE
// @authors     c77cc(c77cc@gmail.com)

package component

const (
	// 微信服务器推送过来的消息类型
	ComponentMsgTypeVerifyTicket = "component_verify_ticket" // 推送 component_verify_ticket 协议
	ComponentMsgTypeUnauthorized = "unauthorized"            // 取消授权的通知
)

type VerifyTicketMessage struct {
	XMLName struct{} `xml:"xml" json:"-"`

	AppId      string `xml:"AppId"      json:"AppId"`
	CreateTime int64  `xml:"CreateTime" json:"CreateTime"`
	InfoType   string `xml:"InfoType"   json:"InfoType"`

	VerifyTicket string `xml:"VerifyTicket" json:"VerifyTicket"`
}

func GetVerifyTicketMessage(msg *MixedMessage) *VerifyTicketMessage {
	return &VerifyTicketMessage{
		AppId:        msg.AppId,
		CreateTime:   msg.CreateTime,
		InfoType:     msg.InfoType,
		VerifyTicket: msg.VerifyTicket,
	}
}

type UnauthorizedMessage struct {
	XMLName struct{} `xml:"xml" json:"-"`

	AppId      string `xml:"AppId"      json:"AppId"`
	CreateTime int64  `xml:"CreateTime" json:"CreateTime"`
	InfoType   string `xml:"InfoType"   json:"InfoType"`

	AuthorizerAppid string `xml:"AuthorizerAppid" json:"AuthorizerAppid"`
}

func GetUnauthorizedMessage(msg *MixedMessage) *UnauthorizedMessage {
	return &UnauthorizedMessage{
		AppId:           msg.AppId,
		CreateTime:      msg.CreateTime,
		InfoType:        msg.InfoType,
		AuthorizerAppid: msg.AuthorizerAppid,
	}
}
