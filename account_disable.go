package yunxin

import (
	"context"
	"fmt"
)

// DisableIMAccountRequest 禁用 IM 账号请求参数结构体
type DisableIMAccountRequest struct {
	Enabled             bool   `json:"enabled"`                         // 是否启用账号：true=启用；false=禁用。禁用后账号无法登录和操作，但仍计入应用账号总数。
	NeedKick            bool   `json:"need_kick,omitempty"`             // 是否强制踢下线：true=踢下线并禁止登录；false=仅禁止登录（默认 false）
	KickNotifyExtension string `json:"kick_notify_extension,omitempty"` // 强制踢下线时的扩展字段（透传至客户端 SDK），最大长度 256 字符
}

type AccountDisableResponse struct {
	AccountId string `json:"account_id,omitempty"` // IM 账号 ID

	Configuration AccountConfiguration `json:"configuration,omitempty"` // 用户账号配置项
}

// AccountDisable 封禁账号
func (c *ImClient) AccountDisable(ctx context.Context, accountId string, param DisableIMAccountRequest) (res Response[AccountDisableResponse], err error) {
	_, _, err = res.Curl(ctx, c, "PATCH", fmt.Sprintf(neteaseBaseURL+"/accounts/%s/actions/disable", accountId), param)

	return
}
