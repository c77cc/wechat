### 获取 jsapi_ticket 示例
```Go
package main

import (
	"fmt"

	"github.com/c77cc/wechat/mp"
	"github.com/c77cc/wechat/mp/jssdk"
)

var AccessTokenServer = mp.NewDefaultAccessTokenServer("appid", "appsecret", nil)
var TicketServer = jssdk.NewDefaultTicketServer(AccessTokenServer, nil)

func main() {
	fmt.Println(TicketServer.Ticket())
}
```