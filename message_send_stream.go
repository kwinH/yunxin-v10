package yunxin

import (
	"context"
	"fmt"
)

// StreamMessage 流式消息主体结构体
type StreamMessageRequest struct {
	Message       StreamMessage  `json:"message"`                  // 消息体
	MessageConfig *MessageConfig `json:"message_config,omitempty"` // 消息相关配置信息（可选）
	RouteConfig   *RouteConfig   `json:"route_config,omitempty"`   // 抄送相关配置信息（可选）
	P2POption     *P2POption     `json:"p2p_option,omitempty"`     // 单聊消息功能相关配置（可选）
	PushConfig    *PushConfig    `json:"push_config,omitempty"`    // 推送相关配置信息（可选）
}

type StreamMessage struct {
	MessageType     int        `json:"message_type"`                // 消息类型，当前仅支持 0（文本消息）
	MessageClientId string     `json:"message_client_id,omitempty"` // 客户端消息 ID（可选）
	MessageServerId int64      `json:"message_server_id,omitempty"` // 首次调用时不填，后续请求需回传
	Attachment      Attachment `json:"attachment"`                  // 流式消息详情
}

// Attachment 流式消息内容结构体
type Attachment struct {
	Text   string `json:"text"`             // 流式消息的内容
	Index  int    `json:"index,omitempty"`  // 流式分片序号（从 0 开始）
	Finish int    `json:"finish,omitempty"` // 是否结束流式输出（0: 未结束，1: 结束）

}

// SendStreamMessage   发送流式消息
// param ownerId          string   发送方账号 ID
// param conversationType int8     会话类型 1：单聊会话
// param to               string   接收方账号 ID
// param param            StreamMessageRequest   发送消息参数
func (c *ImClient) SendStreamMessage(ctx context.Context, ownerId string, conversationType int8, to string, param StreamMessageRequest) (res Response[SendMsgData], err error) {
	conversationId := fmt.Sprintf("%s|%d|%s", ownerId, conversationType, to)

	_, _, err = res.Curl(ctx, c, "POST", fmt.Sprintf("%s/conversations/%s/messages/actions/stream_message", neteaseBaseURL, conversationId), param)

	return
}
