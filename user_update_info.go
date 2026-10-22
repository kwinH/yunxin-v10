package yunxin

import (
	"context"
	"fmt"
)

// UpdateUserInfoRequest 更新用户信息请求参数结构体
type UpdateUserInfoRequest struct {
	Name      string `json:"name,omitempty"`      // 用户昵称，最大长度 64 字符，需通过反垃圾审核
	Avatar    string `json:"avatar,omitempty"`    // 用户头像 URL，最大长度 1024 字符
	Sign      string `json:"sign,omitempty"`      // 用户签名，最大长度 256 字符
	Email     string `json:"email,omitempty"`     // 邮箱地址，格式如 "zhangsan@xx.com"，最大长度 64 字符
	Birthday  string `json:"birthday,omitempty"`  // 用户生日，格式 "xxxx-xx-xx"，最大长度 16 字符
	Mobile    string `json:"mobile,omitempty"`    // 手机号码，非中国大陆需加国家或地区代码（如 +1-xxx 或 +852-xxx）
	Gender    int    `json:"gender,omitempty"`    // 性别：0=未知，1=男，2=女
	Extension string `json:"extension,omitempty"` // 开发者扩展字段，建议封装成 JSON 格式，最大长度 1024 字符

	AntispamConfig AntispamConfiguration `json:"antispam_configuration,omitempty"` // 安全通相关配置项
}

// AntispamConfiguration 安全通配置项
type AntispamConfiguration struct {
	Enabled       bool                `json:"enabled,omitempty"`         // 是否开启安全通，默认 true。仅在已开通安全通业务时生效
	BusinessIdMap []BusinessIdMapping `json:"business_id_map,omitempty"` // 安全通业务 ID 列表
}

// BusinessIdMapping 安全通业务映射
type BusinessIdMapping struct {
	Type       int    `json:"type"`        // 检测类型：1=文本；2=图片
	BusinessId string `json:"business_id"` // 安全通业务 ID
}

type UpdateUserInfoResponse struct {
	AccountId string `json:"account_id,omitempty"` // IM 账号 ID
	UserInfo
}

func (c *ImClient) UpdateUserInfo(ctx context.Context, accountId string, param UpdateUserInfoRequest) (res Response[UpdateUserInfoResponse], err error) {
	_, _, err = res.Curl(ctx, c, "PATCH", fmt.Sprintf("%s/users/%s", neteaseBaseURL, accountId), param)

	return
}
