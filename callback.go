package yunxin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// 事件类型常量定义
const (
	CallbackEventTypeUserInfoUpdate = 3  // 用户资料变更
	CallbackEventTypeP2PMessage     = 1  // 点对点消息
	CallbackEventTypeGroupMessage   = 2  // 群消息
	CallbackEventTypeChatRoom       = 6  // 聊天室消息
	CallbackEventTypeSuperGroup     = 22 // 超大群消息
	CallbackEventTypeRecall         = 35 // 消息撤回
	CallbackEventTypeLogin          = 36 // 登录事件
)

// ChangeOfUserInformation 用户资料变更回调
type ChangeOfUserInformation struct {
	EventType  int    `json:"eventType"`        // 必须 值为3，表示是用户资料变更回调
	Account    string `json:"account"`          // 必须，用户账号
	DeviceID   string `json:"deviceId"`         // 必须，发送方设备id
	ClientType string `json:"clientType"`       // 必须，客户端类型：AOS(1)、IOS(2)、PC(4)、WEB(16)、REST(32)，MAC(64)、HARMONY(65)
	Name       string `json:"name,omitempty"`   // 可选，昵称
	Icon       string `json:"icon,omitempty"`   // 可选，头像
	Sign       string `json:"sign,omitempty"`   // 可选，签名
	Email      string `json:"email,omitempty"`  // 可选，邮箱
	Birth      string `json:"birth,omitempty"`  // 可选，生日
	Mobile     string `json:"mobile,omitempty"` // 可选，手机号
	Gender     int8   `json:"gender,omitempty"` // 可选，性别，0表示未知，1表示男，2表示女
	Ex         string `json:"ex,omitempty"`     // 可选，扩展字段
	Timestamp  string `json:"timestamp"`        // 必须，时间戳
}

// MessageCallback 会话消息回调
type MessageCallback struct {
	EventType      int    `json:"eventType"`                // 回调类型 1-点对点 2-群消息 6-聊天室 22-超大群
	FromAccount    string `json:"fromAccount"`              // 发送方账号
	FromNick       string `json:"fromNick,omitempty"`       // 发送方昵称（可选）
	FromClientType string `json:"fromClientType"`           // 必须，操作者客户端类型：AOS(1)、IOS(2)、PC(4)、WEB(16)、REST(32)，MAC(64)、HARMONY(65)
	FromDeviceID   string `json:"fromDeviceId"`             // 发送设备ID
	To             string `json:"to"`                       // 接收目标：账号/群ID/聊天室ID
	MsgTimestamp   string `json:"msgTimestamp"`             // 消息时间戳
	MsgType        string `json:"msgType"`                  // 消息类型 TEXT/PICTURE等
	FromClientIP   string `json:"fromClientIp,omitempty"`   // 客户端IP（可选）
	FromClientPort string `json:"fromClientPort,omitempty"` // 客户端端口（可选）
	MsgidClient    string `json:"msgidClient,omitempty"`    // 客户端消息ID（可选）
	Body           string `json:"body,omitempty"`           // 消息内容（可选）
	Attach         string `json:"attach,omitempty"`         // 消息附件（可选）
	Ext            string `json:"ext,omitempty"`            // 扩展字段（可选）
}

// RecallCallback 消息撤回回调
type RecallCallback struct {
	EventType      int    `json:"eventType"`                // 必须，固定值35
	FromAccount    string `json:"fromAccount"`              // 必须，操作者账号
	FromDeviceID   string `json:"fromDeviceId"`             // 必须，操作设备ID
	FromClientType string `json:"fromClientType"`           // 必须，操作者客户端类型：AOS(1)、IOS(2)、PC(4)、WEB(16)、REST(32)，MAC(64)、HARMONY(65)
	FromClientIP   string `json:"fromClientIp,omitempty"`   // 可选，客户端IP
	FromClientPort string `json:"fromClientPort,omitempty"` // 可选，客户端端口
	MsgFromAccid   string `json:"msgFromAccid"`             // 必须，消息原发送者
	MsgID          int64  `json:"msgId"`                    // 必须，服务器消息ID
	Time           int64  `json:"time"`                     // 必须，原消息发送时间
	OpeType        int    `json:"opeType"`                  // 必须，7-单聊撤回 8-群撤回
	ToAccount      string `json:"toAccount"`                // 必须，接收方标识
	MsgidClient    string `json:"msgidClient,omitempty"`    // 可选，客户端消息ID
	Msg            string `json:"msg,omitempty"`            // 可选，撤回附言
	Attach         string `json:"attach,omitempty"`         // 可选，扩展字段
	Timestamp      string `json:"timestamp"`                // 必须，操作时间戳（字符串）
}

// LoginCallback 登录回调
type LoginCallback struct {
	EventType        int    `json:"eventType"`                // 必须，固定值36
	FromAccount      string `json:"fromAccount"`              // 必须，操作者账号
	FromDeviceID     string `json:"fromDeviceId"`             // 必须，操作设备ID
	FromClientType   string `json:"fromClientType"`           // 必须，操作者客户端类型：AOS(1)、IOS(2)、PC(4)、WEB(16)、REST(32)，MAC(64)、HARMONY(65)
	FromClientIP     string `json:"fromClientIp,omitempty"`   // 可选，客户端IP
	FromClientPort   string `json:"fromClientPort,omitempty"` // 可选，客户端端口
	Token            string `json:"token"`                    // 必须，登录令牌
	AuthType         int    `json:"authType"`                 // 必须，鉴权方式(0-经典 1-动态 2-第三方)
	LoginExt         string `json:"loginExt"`                 // 必须，登录扩展字段
	CustomTag        string `json:"customTag"`                // 必须，自定义标签
	CustomClientType string `json:"customClientType"`         // 必须，自定义端类型
	AutoLogin        bool   `json:"autoLogin,omitempty"`      // 可选，是否自动登录
	Timestamp        string `json:"timestamp"`                // 必须，操作时间戳
}

type CallbackEvent struct {
	EventType      int    `json:"eventType"`                // 必须
	MsgidClient    string `json:"msgidClient,omitempty"`    // 可选，客户端消息ID
	Attach         string `json:"attach,omitempty"`         // 可选，扩展字段
	FromAccount    string `json:"fromAccount"`              // 必须，操作者账号
	FromDeviceID   string `json:"fromDeviceId"`             // 必须，操作设备ID
	FromClientType string `json:"fromClientType"`           // 必须，操作者客户端类型：AOS(1)、IOS(2)、PC(4)、WEB(16)、REST(32)，MAC(64)、HARMONY(65)
	FromClientIP   string `json:"fromClientIp,omitempty"`   // 可选，客户端IP
	FromClientPort string `json:"fromClientPort,omitempty"` // 可选，客户端端口
	Timestamp      string `json:"timestamp"`                // 必须，操作时间戳

	//用户资料变更回调
	*ChangeOfUserInformation
	//会话消息回调
	*MessageCallback
	//消息撤回回调
	*RecallCallback
	//登录回调
	*LoginCallback
}

type CallbackResponse struct {
	// 0: 允许执行 1: 拒绝执行
	ErrCode int `json:"errCode"`

	// 当 ErrCode=1 时有效，范围 [20000, 20099]
	// 对于消息类型的第三方回调（eventType=1、2、6、22、41、72、73、74），支持设置为 200 的错误码，客户端表现为消息发送成功，其实消息发送失败
	ResponseCode int `json:"responseCode"`

	// 对于消息类型的第三方回调有效（eventType=1、2、6、22、41、72、73、74），用于修改消息内容 (JSON 格式)
	ModifyResponse ModifyPayload `json:"modifyResponse"`

	// 对于消息类型的第三方回调有效（eventType=1、2、6、22），用于传递第三方回调的扩展信息，最大 1,024 个字符
	CallbackExt string `json:"callbackExt"`
}

type ModifyPayload struct {
	// 通用消息字段
	Body   string `json:"body,omitempty"`   // 消息内容
	Attach string `json:"attach,omitempty"` // 附件信息
	Ext    string `json:"ext,omitempty"`    // 扩展字段

	// 推送相关字段
	PushContent string `json:"pushContent,omitempty"` // 推送文案
	PushPayload string `json:"pushPayload,omitempty"` // 推送负载

	// 功能控制字段
	PushEnable    *bool `json:"pushEnable,omitempty"`    // 是否推送
	NeedPushNick  *bool `json:"needPushNick,omitempty"`  // 推送昵称
	PersistEnable *bool `json:"persistEnable,omitempty"` // 是否持久化

	// 聊天室特殊字段
	SkipHistory   *bool `json:"skipHistory,omitempty"`   // 跳过历史记录(仅聊天室)
	RoamingEnable *bool `json:"roamingEnable,omitempty"` // 是否支持漫游
	HistoryEnable *bool `json:"historyEnable,omitempty"` // 是否存历史

	// 圈组系统消息字段
	Msg         string `json:"msg,omitempty"`         // 系统消息内容
	OperatorExt string `json:"operatorExt,omitempty"` // 操作者扩展信息

	// 聊天室标签配置
	Tag             []string `json:"tag,omitempty"`             // 用户标签
	NotifyTargetTag string   `json:"notifyTargetTag,omitempty"` // 标签表达式
}

func (r *CallbackEvent) GetExt() map[string]any {
	ext := make(map[string]any)

	// 根据事件类型获取对应的扩展字段
	var extStr string
	switch r.EventType {
	case CallbackEventTypeP2PMessage, CallbackEventTypeGroupMessage,
		CallbackEventTypeChatRoom, CallbackEventTypeSuperGroup:
		if r.MessageCallback != nil {
			extStr = r.MessageCallback.Ext
		}
	case CallbackEventTypeRecall:
		if r.RecallCallback != nil {
			extStr = r.RecallCallback.Attach // 撤回回调的扩展字段在attach
		}
	case CallbackEventTypeUserInfoUpdate:
		if r.ChangeOfUserInformation != nil {
			extStr = r.ChangeOfUserInformation.Ex
		}
	}

	if extStr != "" {
		json.Unmarshal([]byte(extStr), &ext)
	}

	return ext
}

func (r *CallbackEvent) TimestampToTime() *time.Time {
	// 将字符串时间戳转换为 int64
	timestamp, err := strconv.ParseInt(r.Timestamp, 10, 64)
	if err != nil {
		return nil // 或返回默认时间
	}

	// 处理毫秒级时间戳
	sec := timestamp / 1000
	msec := timestamp % 1000

	// 创建 UTC 时间对象
	utcTime := time.Unix(sec, msec*int64(time.Millisecond)).UTC()

	// 转换为北京时间 (UTC+8)
	cst := time.FixedZone("CST", 8*60*60)
	beijingTime := utcTime.In(cst)

	return &beijingTime
}

func (r *CallbackEvent) UnmarshalJSON(data []byte) error {
	// 使用中间类型解除递归
	type Alias CallbackEvent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	// 先解析基础字段
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// 根据事件类型初始化具体回调结构
	switch r.EventType {
	case CallbackEventTypeUserInfoUpdate:
		if r.ChangeOfUserInformation == nil {
			r.ChangeOfUserInformation = &ChangeOfUserInformation{}
		}
		return json.Unmarshal(data, r.ChangeOfUserInformation)
	case CallbackEventTypeP2PMessage, CallbackEventTypeGroupMessage, CallbackEventTypeChatRoom, CallbackEventTypeSuperGroup:
		if r.MessageCallback == nil {
			r.MessageCallback = &MessageCallback{}
		}
		err := json.Unmarshal(data, r.MessageCallback)
		r.Timestamp = r.MsgTimestamp
		return err
	case CallbackEventTypeRecall:
		if r.RecallCallback == nil {
			r.RecallCallback = &RecallCallback{}
		}
		return json.Unmarshal(data, r.RecallCallback)
	case CallbackEventTypeLogin:
		if r.LoginCallback == nil {
			r.LoginCallback = &LoginCallback{}
		}
		return json.Unmarshal(data, r.LoginCallback)
	}

	return nil
}

func (c *ImClient) Callback(req *http.Request, fun func(param *CallbackEvent) (CallbackResponse, error)) (res CallbackResponse, err error) {
	bodyBytes, err := c.CheckSumMd5(req)
	if err != nil {
		res = CallbackResponse{
			ErrCode:     1,
			CallbackExt: "回调：验签失败",
		}
		return
	}

	param := new(CallbackEvent)
	err = json.Unmarshal(bodyBytes, param)
	if err != nil {
		res = CallbackResponse{
			ErrCode:     1,
			CallbackExt: "参数错误",
		}
		return
	}

	return fun(param)
}
