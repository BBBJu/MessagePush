package service

import (
	"context"
	"fmt"
	"messagePush/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var MyLarkSender *LarkSender

func init() {
	MyLarkSender = NewLarkSender(config.MyConfig.AppId, config.MyConfig.AppSecret)
}

func NewLarkSender(appId string, appSecret string) *LarkSender {
	client := lark.NewClient(appId, appSecret)
	return &LarkSender{
		client: client,
	}

	//utils.GetTenantAccessToken(client)
	// t := service.MessageParams{
	// 	ReceiveId: utils.ReceiveId,
	// 	Content:   "{\"text\":\"不会哈气学哈气，机密咋摆你咋摆233\"}", // 转义双引号
	// 	MsgType:   "text",
	// }
}

type LarkSender struct {
	client *lark.Client
}

func (l *LarkSender) SendMessage(param MessageParams) error {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`open_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(param.ReceiveId).
			MsgType(param.MsgType).
			Content(param.Content).
			Build()).
		Build()

	// 发起请求
	resp, err := l.client.Im.V1.Message.Create(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return err
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
	return nil
}
