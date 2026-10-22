package yunxin

import (
	"context"
)

// CreateIMAccountRequest 创建 IM 账号请求参数结构体
type CreateIMAccountRequest struct {
	AccountId string `json:"account_id"`      // IM 账号，请确保唯一性。只允许字母、数字、半角下划线、@、半角点以及半角横线。长度上限为 32 字符。
	Token     string `json:"token,omitempty"` // IM 账号对应的登录密钥 Token。未指定则自动生成，永久有效。长度上限 128 字符。

	Configuration AccountConfiguration `json:"configuration,omitempty"`    // IM 账号配置项
	UserInfo      UserInfo             `json:"user_information,omitempty"` // 用户信息
}

// AccountConfiguration IM 账号配置项
type AccountConfiguration struct {
	Enabled                      *bool `json:"enabled,omitempty"`                          // 账号是否可用，默认 true
	P2PChatBanned                bool  `json:"p2p_chat_banned,omitempty"`                  // 单聊禁言标记：true=禁言
	TeamChatBanned               bool  `json:"team_chat_banned,omitempty"`                 // 群聊禁言标记：true=禁言
	ChatroomChatBanned           bool  `json:"chatroom_chat_banned,omitempty"`             // 聊天室禁言标记：true=禁言
	QChatChatBanned              bool  `json:"qchat_chat_banned,omitempty"`                // 圈组禁言标记：true=禁言
	PushEnabledWhenDesktopOnline *bool `json:"push_enabled_when_desktop_online,omitempty"` // 桌面端在线时是否允许移动端推送，默认 true。仅当设备有登录记录时才支持设置
}

// UserInfo 用户信息结构体
type UserInfo struct {
	Name               string `json:"name,omitempty"`                 // 用户昵称，最大 64 字符，需通过反垃圾审核
	Avatar             string `json:"avatar,omitempty"`               // 用户头像 URL，最大 1024 字符
	Sign               string `json:"sign,omitempty"`                 // 用户签名，最大 256 字符
	Email              string `json:"email,omitempty"`                // 邮箱地址，格式如 "zhangsan@xx.com"
	Birthday           string `json:"birthday,omitempty"`             // 生日，格式 "xxxx-xx-xx"，最大 16 字符
	Mobile             string `json:"mobile,omitempty"`               // 手机号码，非中国大陆需加国家代码，如 +1-xxx 或 +852-xxx
	Gender             int    `json:"gender,omitempty"`               // 性别：0=未知，1=男，2=女
	Extension          string `json:"extension,omitempty"`            // 开发者扩展字段，建议 JSON 格式，最大 1024 字符
	AntispamBusinessId string `json:"antispam_business_id,omitempty"` // 安全通业务 ID，JSON 字符串 {"textbid":"","picbid":""}
}

// AccountResponse 更新 IM 账号请求参数结构体
type AccountResponse struct {
	AccountId string `json:"account_id,omitempty"` // IM 账号 ID

	Token string `json:"token,omitempty"` // IM 账号对应的登录密钥

	Configuration AccountConfiguration `json:"configuration,omitempty"`    // 用户账号配置项
	UserInfo      UserInfo             `json:"user_information,omitempty"` // 用户信息
}

// CreateAccount 注册 IM 账号
func (c *ImClient) CreateAccount(ctx context.Context, param CreateIMAccountRequest) (res Response[AccountResponse], err error) {
	_, _, err = res.Curl(ctx, c, "POST", neteaseBaseURL+"/accounts", param)

	return
}
