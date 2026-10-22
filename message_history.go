package yunxin

import (
	"context"
	"fmt"
)

// QueryMessagesRequest 查询消息请求参数结构体
type QueryMessagesRequest struct {
	BeginTime         int64  `url:"begin_time"`                     // 查询开始时间（毫秒）
	EndTime           int64  `url:"end_time"`                       // 查询截止时间（毫秒）
	PageToken         string `url:"page_token,omitempty"`           // 分页标识符，如果为空，则从第一页开始查询
	Limit             int    `url:"limit"`                          // 每页返回的消息数上限，最大为 100
	Descending        bool   `url:"descending,omitempty"`           // 是否按时间正序：true=正序；false=倒序
	MessageType       string `url:"message_type,omitempty"`         // 指定需要查询的消息类型，多个用逗号分隔。例如："0,1,2"，其中支持的消息类型如下：0：文本；1：图片；2：语音；3：视频；4：地理位置；5：通知；6：文件；10：提示；11：Robot；100：自定义。
	CheckTeamValid    bool   `url:"check_team_valid,omitempty"`     // 是否检查高级群成员有效性，默认 true
	IncludeNoSenseMsg bool   `url:"include_no_sense_msg,omitempty"` // 查询结果是否包含无感知消息，默认 false
}

// QueryMessagesResponse 查询消息返回结果结构体
type QueryMessagesResponse struct {
	HasMore   bool   `json:"has_more"`             // 是否还有更多数据
	NextToken string `json:"next_token,omitempty"` // 分页标识符（可选）

	Items []MessageItem `json:"items"` // 返回的历史消息列表
}

// MessageItem 消息项结构体
type MessageItem struct {
	MessageServerId int64  `json:"message_server_id"` // 服务端消息 ID
	SenderId        string `json:"sender_id"`         // 消息发送方账号 ID
	MessageType     int    `json:"message_type"`      // 消息类型：0=文本；1=图片；2=语音；3=视频；4=地址位置；5=通知；6=文件；10=提示；11=机器人；100=自定义
	CreateTime      int64  `json:"create_time"`       // 消息发送时间戳（毫秒）
	MessageClientId string `json:"message_client_id"` // 客户端消息 ID

	SenderClientType string `json:"sender_client_type,omitempty"` // 发送端类型：1=Android；2=iOS；4=PC；16=Web；32=REST；64=MAC
	Text             string `json:"text,omitempty"`               // 文本内容或多媒体描述文本
	Attachment       any    `json:"attachment,omitempty"`         // 多媒体属性或自定义消息内容（可根据实际类型使用 map 或 struct）
	Extension        string `json:"extension,omitempty"`          // 第三方扩展字段
	TeamId           int64  `json:"team_id,omitempty"`            // 群组 ID（与 receiver_id 只会返回一个）
	ReceiverId       string `json:"receiver_id,omitempty"`        // 消息接收者 ID

	SubType         int    `json:"sub_type,omitempty"`          // 自定义消息子类型（仅 message_type = 100 时有效）
	ModifyAccountId string `json:"modify_account_id,omitempty"` // 消息更新者的用户账号 ID（消息被更新时才返回）
	ModifyTime      string `json:"modify_time,omitempty"`       // 消息更新时间（消息被更新时才返回）

	ThreadConfig *ThreadConfig  `json:"thread_config,omitempty"` // Thread 消息配置项（可选）
	ThreadRoot   *ThreadMessage `json:"thread_root,omitempty"`   // 根消息对象（可选）
	ThreadReply  *ThreadMessage `json:"thread_reply,omitempty"`  // 被回复的消息对象（可选）
}

// GetMessagesPage 分页查询历史消息
func (c *ImClient) GetMessagesPage(ctx context.Context, ownerId string, conversationType int8, to string, param QueryMessagesRequest) (res Response[QueryMessagesResponse], err error) {
	conversationId := fmt.Sprintf("%s|%d|%s", ownerId, conversationType, to)

	// 构建带参数的 URL
	userAddr := fmt.Sprintf("%s/conversations/%s/messages", neteaseBaseURL, conversationId)

	_, _, err = res.Curl(ctx, c, "GET", userAddr, param)

	return
}
