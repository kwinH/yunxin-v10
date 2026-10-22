package yunxin

import (
	"context"
)

// ModifyMsgRequest 定义更新消息的请求参数结构体
type ModifyMsgRequest struct {
	Operator  string        `json:"operator"`            // 更新消息的操作者账号 ID（必填）
	Extension string        `json:"extension,omitempty"` // 开发者扩展字段，JSON 格式（可选）
	Type      int           `json:"type"`                // 原始会话类型（1:单聊, 2:高级群, 3:超大群）（必填）
	Message   UpdateMessage `json:"message"`             // 消息体（必填）
}

// UpdateMessage 消息主体结构，用于更新消息内容
type UpdateMessage struct {
	MessageServerId int64                `json:"message_server_id"`          // 待更新消息的原始服务器 ID（必填）
	ReceiverId      string               `json:"receiver_id"`                // 原始接收者 ID 或群组 ID（必填）
	Time            int64                `json:"time"`                       // 原始发送时间（毫秒）（必填）
	SenderId        string               `json:"sender_id"`                  // 原始消息发送者账号 ID（必填）
	MessageType     int                  `json:"message_type"`               // 原始消息类型（0-文本, 1-图片等）（必填）
	SubType         int                  `json:"sub_type,omitempty"`         // 自定义消息子类型（仅 message_type=100 时有效）
	Text            string               `json:"text,omitempty"`             // 新文本内容或描述信息
	Attachment      any                  `json:"attachment,omitempty"`       // 非文本消息属性或自定义内容
	MessageConfig   *ModifyMsgConfig     `json:"message_config,omitempty"`   // 消息配置项
	RouteConfig     *RouteConfig         `json:"route_config,omitempty"`     // 抄送相关配置项
	PushConfig      *PushConfig          `json:"push_config,omitempty"`      // 推送相关配置项
	AntispamConfig  *AntispamConfig      `json:"antispam_config,omitempty"`  // 安全通相关配置项
	P2POption       *P2POption           `json:"p2p_option,omitempty"`       // 单聊消息功能配置项
	TeamOption      *ModifyMsgTeamOption `json:"team_option,omitempty"`      // 高级群消息配置项
	SuperteamOption *SuperteamOption     `json:"superteam_option,omitempty"` // 超大群消息配置项
}

type ModifyMsgConfig struct {
	ConversationUpdateEnabled bool `json:"conversation_update_enabled,omitempty"` // 是否更新至会话列表
}

type ModifyMsgTeamOption struct {
	IgnoreChatBanned       bool `json:"ignore_chat_banned,omitempty"`        // 忽略群禁言
	IgnoreMemberChatBanned bool `json:"ignore_member_chat_banned,omitempty"` // 忽略成员禁言
	CheckTeamMemberValid   bool `json:"check_team_member_valid,omitempty"`   // 验证群成员身份
}

//------------响应------------------

// ModifyMsgData 响应中的 data 字段结构
type ModifyMsgData struct {
	ModifyAccountId string `json:"modify_account_id,omitempty"` // 消息更新者的用户账号 ID（仅在消息被更新时返回）
	ModifyTime      int64  `json:"modify_time,omitempty"`       // 消息更新时间戳（毫秒）
	MessageServerId int64  `json:"message_server_id,omitempty"` // 服务端消息 ID
	SenderId        string `json:"sender_id,omitempty"`         // 消息发送方账号 ID
	ReceiverId      string `json:"receiver_id,omitempty"`       // 接收者账号 ID
	CreateTime      int64  `json:"create_time,omitempty"`       // 消息发送时间戳（毫秒）
	MessageType     int    `json:"message_type"`                // 消息类型（必返）
	SubType         int    `json:"sub_type,omitempty"`          // 自定义消息子类型
	Text            string `json:"text,omitempty"`              // 文本/提示内容或描述文本
	Attachment      any    `json:"attachment,omitempty"`        // 多媒体属性或自定义消息内容
}

// ModifyMsg 修改消息
// param: 修改消息的参数结构体
func (c *ImClient) ModifyMsg(ctx context.Context, param ModifyMsgRequest) (res Response[ModifyMsgData], err error) {
	_, _, err = res.Curl(ctx, c, "POST", neteaseBaseURL+"/messages/actions/modifyMsg", param)

	return
}
