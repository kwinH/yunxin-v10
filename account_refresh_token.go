package yunxin

import (
	"context"
	"fmt"
)

type RefreshTokenResponse struct {
	AccountId string `json:"account_id,omitempty"` // IM 账号 ID
	Token     string `json:"token,omitempty"`      // IM 账号对应的登录密钥

	Configuration AccountConfiguration `json:"configuration,omitempty"` // 用户账号配置项
}

// RefreshToken 刷新 Token
func (c *ImClient) RefreshToken(ctx context.Context, accountId string) (res Response[RefreshTokenResponse], err error) {
	_, _, err = res.Curl(ctx, c, "PATCH", fmt.Sprintf(neteaseBaseURL+"/accounts/%s/actions/refresh_token", accountId), nil)

	return
}
