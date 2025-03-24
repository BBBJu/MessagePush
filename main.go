package main

import (
	"messagePush/config"
	"messagePush/service"
	"messagePush/utils"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

func main() {
	// 创建 Client
	AppID, AppSecret := config.MyConfig.AppId, config.MyConfig.AppSecret
	client := lark.NewClient(AppID, AppSecret)
	utils.GetTenantAccessToken(client)
	t := service.MessageParams{
		ReceiveId: utils.ReceiveId,
		Content:   "{\"text\":\"不会哈气学哈气，机密咋摆你咋摆\"}", // 转义双引号
		MsgType:   "text",
	}
	service.SendMessage(client, t)
}
