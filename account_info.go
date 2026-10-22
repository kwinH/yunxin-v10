package yunxin

import (
	"context"
	"fmt"
)

// AccountInfoResponse 查询 IM 账号信息返回结构体
type AccountInfoResponse struct {
	AccountId     string               `json:"account_id,omitempty"`    // IM 账号 ID
	Configuration AccountConfiguration `json:"configuration,omitempty"` // 用户账号配置项
}

// AccountInfo 获取 IM 账号信息
func (c *ImClient) AccountInfo(ctx context.Context, accountId string) (res Response[AccountInfoResponse], err error) {
	_, _, err = res.Curl(ctx, c, "GET", fmt.Sprintf("%s/accounts/%s", neteaseBaseURL, accountId), nil)

	return
}
