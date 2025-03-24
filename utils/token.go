package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"messagePush/config"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauth "github.com/larksuite/oapi-sdk-go/v3/service/auth/v3"
)

type TokenResponse struct {
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func GetTenantAccessToken(client *lark.Client) (string, error) {
	AppID, AppSecret := config.MyConfig.AppId, config.MyConfig.AppSecret
	req := larkauth.NewInternalTenantAccessTokenReqBuilder().
		Body(larkauth.NewInternalTenantAccessTokenReqBodyBuilder().
			AppId(AppID).
			AppSecret(AppSecret).
			Build()).
		Build()

	// 发起请求
	resp, err := client.Auth.V3.TenantAccessToken.Internal(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return "", err
	}
	// 业务处理
	fmt.Println(larkcore.Prettify(resp))

	// 定义响应结构体

	// 解析RawBody
	var tokenResp TokenResponse
	if err := json.Unmarshal(resp.RawBody, &tokenResp); err != nil {
		fmt.Printf("JSON解析失败: %v", err)
		return "", err
	}
	fmt.Printf("访问令牌: %s\n有效期: %d秒", tokenResp.TenantAccessToken, tokenResp.Expire)
	return tokenResp.TenantAccessToken, nil
}
