package yunxin

import (
	"context"
	"encoding/json"
	"fmt"
)

// SystemUser 系统用户账号 ID
const SystemUser = "system"

// NotificationConfig 通知配置项
type NotificationConfig struct {
	OfflineEnabled *bool `json:"offline_enabled,omitempty"` // 是否需要存离线消息，默认为 true（存离线）
	UnreadEnabled  *bool `json:"unread_enabled,omitempty"`  // 是否计入未读数，默认为 true（计入）
}

// CustomNotificationRequest 定义发送自定义系统通知的请求参数结构体
type CustomNotificationRequest struct {
	SenderId           string              `json:"sender_id"`                     // 发送者 IM 账号 ID（必填）
	Type               int                 `json:"type"`                          // 自定义系统通知类型。1：单聊系统通知。2：高级群系统通知。3：超大群系统通知（必填）
	ReceiverId         string              `json:"receiver_id"`                   // 系统通知接收者 ID。type=1 时该参数为接收者账号 ID；type=2 或 3 时该参数为接收的群 ID（必填）
	Content            string              `json:"content"`                       // 自定义系统通知的内容，由开发者自行组装的 JSON 格式字符串，长度上限 4096 位字符（必填）
	Sound              string              `json:"sound,omitempty"`               // 指定客户端本地声音文件名，长度上限 30 位字符
	NotificationConfig *NotificationConfig `json:"notification_config,omitempty"` // 通知配置项，不填则采用默认值
	PushConfig         *PushConfig         `json:"push_config,omitempty"`         // 推送相关配置项
	RouteConfig        *RouteConfig        `json:"route_config,omitempty"`        // 抄送相关配置项
}

// SendCustomNotification 发送自定义系统通知
// param param   CustomNotificationRequest   发送自定义系统通知参数
func (c *ImClient) SendCustomNotification(ctx context.Context, param CustomNotificationRequest) (res Response[struct{}], err error) {
	_, _, err = res.Curl(ctx, c, "POST", neteaseBaseURL+"/custom_notification", param)

	return
}

func (c *ImClient) SendSystemNotification(ctx context.Context, receiverId string, code int, msg string, data any) (err error) {
	payload := map[string]any{
		"code": code,
		"msg":  msg,
	}
	if data != nil {
		payload["data"] = data
	}
	content, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("构造云信系统通知内容失败:%+v", err)
	}
	contentStr := string(content)
	res, err := c.SendCustomNotification(ctx, CustomNotificationRequest{
		Type:       1,
		SenderId:   SystemUser,
		ReceiverId: receiverId,
		Content:    contentStr,
	})

	if err != nil {
		return fmt.Errorf("发送云信系统通知失败:%+v", err)
	}

	if res.Code != 200 {
		return fmt.Errorf("发送云信系统通知失败:%+v", res)
	}

	return
}
