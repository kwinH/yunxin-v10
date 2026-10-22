package yunxin

import (
	"context"
	"fmt"
)

// SendMessageRequest 定义发送消息的请求参数结构体
type SendMessageRequest struct {
	Message         Message          `json:"message"`                    // 消息体
	MessageConfig   *MessageConfig   `json:"message_config,omitempty"`   // 消息配置项
	TeamOption      *TeamOption      `json:"team_option,omitempty"`      // 高级群消息功能配置项
	SuperteamOption *SuperteamOption `json:"superteam_option,omitempty"` // 超大群消息功能配置项
	RobotConfig     *RobotConfig     `json:"robot_config,omitempty"`     // 机器人功能配置项
	AiParams        *AiParams        `json:"ai_params,omitempty"`        // 数字人功能配置项
	ThreadConfig    *ThreadConfig    `json:"thread_config,omitempty"`    // Thread 消息配置项
	Extension       string           `json:"extension,omitempty"`        // 开发者扩展字段，JSON 格式
}

// Message 消息主体结构
type Message struct {
	MessageType int         `json:"message_type"`                // 消息类型：0-文本，1-图片，2-语音，3-视频，4-地理位置，6-文件，10-提示，100-自定义
	SubType     int         `json:"sub_type,omitempty"`          // 自定义消息子类型，仅 message_type=100 时有效
	Text        string      `json:"text,omitempty"`              // 文本内容或描述信息
	Attachment  interface{} `json:"attachment,omitempty"`        // 非文本消息属性或自定义内容
	ClientId    string      `json:"message_client_id,omitempty"` // 客户端消息 ID
}

// MessageConfig 消息配置项
type MessageConfig struct {
	UnreadEnabled             bool `json:"unread_enabled,omitempty"`              // 是否计入未读数
	MutilSyncEnabled          bool `json:"mutil_sync_enabled,omitempty"`          // 是否多端同步
	OfflineEnabled            bool `json:"offline_enabled,omitempty"`             // 是否存离线消息
	HistoryEnabled            bool `json:"history_enabled,omitempty"`             // 是否存云端历史
	RoamingEnabled            bool `json:"roaming_enabled,omitempty"`             // 是否漫游
	ConversationUpdateEnabled bool `json:"conversation_update_enabled,omitempty"` // 是否更新至会话列表
}

// TargetOption 群定向消息配置
type TargetOption struct {
	ReceiverAccountIds   []string `json:"receiver_account_ids"`              // 群定向消息成员列表，即指定接收群消息的群成员列表。
	Inclusive            *bool    `json:"inclusive,omitempty"`               // 群定向消息成员列表是否为可见列表。默认为 true，即 receiver_account_ids 为可见（接收）列表。
	CheckTeamMemberValid *bool    `json:"check_team_member_valid,omitempty"` // 发送群定向消息时，是否需要验证定向成员的群成员身份。默认为 true，即需要验证群成员身份。
	VisibleToNewMember   bool     `json:"visible_to_new_member,omitempty"`   // 新进群成员是否可见该消息。默认为 false，即新进群成员不可以查看该定向消息。若设置为 true，则新进群成员若可以查询该定向消息，可以通过云端历史相关接口查询到该消息。
}

// RouteConfig 抄送相关配置
type RouteConfig struct {
	RouteEnabled     bool   `json:"route_enabled,omitempty"`     // 是否抄送
	RouteEnvironment string `json:"route_environment,omitempty"` // 抄送环境名称
}

// PushConfig 推送相关配置
type PushConfig struct {
	PushEnabled          bool     `json:"push_enabled,omitempty"`           // 是否推送
	PushNickEnabled      bool     `json:"push_nick_enabled,omitempty"`      // 推送是否带昵称
	PushContent          string   `json:"push_content,omitempty"`           // 自定义推送文案
	PushPayload          string   `json:"push_payload,omitempty"`           // 推送 payload
	PushForcePushEnabled bool     `json:"push_forcepush_enable,omitempty"`  // 是否强制推送
	PushForcePushAll     bool     `json:"push_forcepush_all,omitempty"`     // 强制推送所有人
	PushForcePushIds     []string `json:"push_forcepush_ids,omitempty"`     // 强制推送指定用户
	PushForcePushContent string   `json:"push_forcepush_content,omitempty"` // 强制推送文案
}

// AntispamConfig 安全通相关配置
type AntispamConfig struct {
	AntispamEnabled              bool   `json:"antispam_enabled,omitempty"`                // 是否过审核
	AntispamBusinessId           string `json:"antispam_business_id,omitempty"`            // 安全通业务 ID
	AntispamExtension            string `json:"antispam_extension,omitempty"`              // 反垃圾增强版参数
	AntispamCheating             string `json:"antispam_cheating,omitempty"`               // 反作弊参数
	AntispamCustomMessageEnabled bool   `json:"antispam_custom_message_enabled,omitempty"` // 自定义消息是否审核
	AntispamCustomMessage        string `json:"antispam_custom_message,omitempty"`         // 自定义审核内容
}

// P2POption 单聊消息功能配置
type P2POption struct {
	CheckFriend bool `json:"check_friend,omitempty"` // 是否只发给好友
}

// TeamOption 高级群消息配置
type TeamOption struct {
	MarkAsRead             bool `json:"mark_as_read,omitempty"`              // 是否启用已读功能
	IgnoreChatBanned       bool `json:"ignore_chat_banned,omitempty"`        // 忽略群禁言
	IgnoreMemberChatBanned bool `json:"ignore_member_chat_banned,omitempty"` // 忽略成员禁言
	CheckTeamMemberValid   bool `json:"check_team_member_valid,omitempty"`   // 验证群成员身份
}

// SuperteamOption 超大群消息配置
type SuperteamOption struct {
	IgnoreChatBanned       bool `json:"ignore_chat_banned,omitempty"`        // 忽略超大群禁言
	IgnoreMemberChatBanned bool `json:"ignore_member_chat_banned,omitempty"` // 忽略成员禁言
	CheckTeamMemberValid   bool `json:"check_team_member_valid,omitempty"`   // 验证群成员身份
}

// RobotConfig 机器人配置
type RobotConfig struct {
	RobotAccountId     string `json:"robot_account_id,omitempty"`     // 机器人账号
	RobotTopic         string `json:"robot_topic,omitempty"`          // 机器人话题
	RobotFunction      string `json:"robot_function,omitempty"`       // 机器人功能
	RobotCustomContent string `json:"robot_custom_content,omitempty"` // 自定义内容
}

// AiParams 数字人功能配置
type AiParams struct {
	Account         string                 `json:"account"`                    // 数字人账号
	Content         *AiContent             `json:"content,omitempty"`          // 查询内容
	Messages        []*AiMessage           `json:"messages,omitempty"`         // 回溯消息列表
	PromptVariables interface{}            `json:"prompt_variables,omitempty"` // 提示词变量
	Config          map[string]interface{} `json:"config,omitempty"`           // 模型配置参数
}

// AiContent 数字人查询内容
type AiContent struct {
	Msg  string `json:"msg"`  // 内容
	Type int    `json:"type"` // 类型（当前仅支持 0）
}

// AiMessage 回溯消息对象
type AiMessage struct {
	Msg  string `json:"msg"`  // 消息内容
	Type int    `json:"type"` // 类型
	Role string `json:"role"` // 角色（user/assistant）
}

// ThreadConfig Thread 消息配置
type ThreadConfig struct {
	ThreadRoot  *ThreadMessage `json:"thread_root,omitempty"`  // 根消息
	ThreadReply *ThreadMessage `json:"thread_reply,omitempty"` // 被回复的消息
}

// ThreadMessage Thread 消息对象
type ThreadMessage struct {
	SenderId        string `json:"sender_id,omitempty"`         // 发送者 ID
	ReceiverId      string `json:"receiver_id,omitempty"`       // 接收者 ID
	CreateTime      int64  `json:"create_time,omitempty"`       // 创建时间（毫秒）
	MessageServerId int64  `json:"message_server_id,omitempty"` // 服务器消息 ID
	MessageClientId string `json:"message_client_id,omitempty"` // 客户端消息 ID
}

//------------响应------------------

// SendMsgData 响应中的 data 字段结构
type SendMsgData struct {
	MessageServerId  int64         `json:"message_server_id,omitempty"` // 服务端消息 ID
	MessageClientId  string        `json:"message_client_id"`           // 客户端消息 ID（必返）
	SenderId         string        `json:"sender_id,omitempty"`         // 消息发送方账号 ID
	ConversationType int           `json:"conversation_type,omitempty"` // 会话类型
	ReceiverId       string        `json:"receiver_id,omitempty"`       // 接收者账号 ID
	CreateTime       int64         `json:"create_time,omitempty"`       // 消息发送时间戳（毫秒）
	MessageType      int           `json:"message_type"`                // 消息类型（必返）
	SubType          int           `json:"sub_type,omitempty"`          // 自定义消息子类型
	Text             string        `json:"text,omitempty"`              // 文本/提示内容或描述文本
	Attachment       interface{}   `json:"attachment,omitempty"`        // 多媒体属性或自定义消息内容
	ThreadConfig     *ThreadConfig `json:"thread_config,omitempty"`     // Thread 消息配置项
}

// SendMessage  发送消息
// param ownerId          string   发送方账号 ID
// param conversationType int8     会话类型 1：单聊会话。2：高级群会话。3：超大群会话
// param to               string   接收方账号 ID
// param param            SendMessageRequest   发送消息参数
func (c *ImClient) SendMessage(ctx context.Context, ownerId string, conversationType int8, to string, param SendMessageRequest) (res Response[SendMsgData], err error) {
	conversationId := fmt.Sprintf("%s|%d|%s", ownerId, conversationType, to)

	_, _, err = res.Curl(ctx, c, "POST", fmt.Sprintf("%s/conversations/%s/messages", neteaseBaseURL, conversationId), param)

	return
}
