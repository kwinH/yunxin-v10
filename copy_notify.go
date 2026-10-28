package yunxin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	//EventTypeConversation 表示CONVERSATION消息，即会话类型的消息（目前包括P2P聊天消息，群组聊天消息，群组操作，好友操作）
	EventTypeConversation = "1"
	//EventTypeLogin 表示LOGIN消息，即用户登录事件的消息
	EventTypeLogin = "2"
	//EventTypeLogout  表示LOGOUT消息，即用户登出事件的消息
	EventTypeLogout = "3"
	//EventTypeChatRoom 表示CHATROOM消息，即聊天室中聊天的消息
	EventTypeChatRoom = "4"
	//EventTypeMediaDuration 汇报实时音视频通话时长、白板事件时长的消息
	EventTypeMediaDuration = "5"
	//EventTypeMediaInfo 汇报音视频/白板文件的大小、下载地址等消息
	EventTypeMediaInfo = "6"
	//EventTypeP2PMessageRecall 单聊消息撤回抄送
	EventTypeP2PMessageRecall = "7"
	//EventTypeGroupMessageRecall 群聊消息撤回抄送
	EventTypeGroupMessageRecall = "8"
	//EventTypeChatRoomInOut 汇报主播或管理员进出聊天室事件消息
	EventTypeChatRoomInOut = "9"
	//EventTypeECPCallback 汇报专线电话通话结束回调抄送的消息
	EventTypeECPCallback = "10"
	//EventTypeSMSCallback 汇报短信回执抄送的消息
	EventTypeSMSCallback = "11"
	//EventTypeSMSReply 汇报短信上行消息
	EventTypeSMSReply = "12"
	//EventTypeAvRoomInOut 汇报用户进出音视频/白板房间的消息
	EventTypeAvRoomInOut = "13"
	//EventTypeChatRoomQueueOperate 汇报聊天室队列操作的事件消息
	EventTypeChatRoomQueueOperate = "14"
)

// LoginEventCopyInfo 登录事件消息抄送
type LoginEventCopyInfo struct {
	EventType  string `json:"eventType"`  //值为2，表示是登录事件的消息
	AcctID     string `json:"accid"`      //发生登录事件的用户帐号，字符串类型
	IPAdrees   string `json:"clientIp"`   //登录时的ip地址
	ClientType string `json:"clientType"` //客户端类型： AOS、IOS、PC、WINPHONE、WEB、REST，字符串类型
	Code       string `json:"code"`       //登录事件的返回码，可转为Integer类型数据
	SdkVersion string `json:"sdkVersion"` //当前sdk的版本信息，字符串类型
	Timestamp  string `json:"timestamp"`  //登录事件发生时的时间戳，可转为Long型数据
}

// LogoutEventCopyInfo 登出事件消息抄送
type LogoutEventCopyInfo struct {
	LoginEventCopyInfo      //eventType值为3，表示是登出事件的消息
	LogOutReason       int8 `json:"logOutReason"` //登出原因： 1：用户注销 2：用户断开连接 3：用户被自己其它端踢下线 4：根据互踢策略被踢下线
}

// MsgCopyInfo 会话类型信息抄送
type MsgCopyInfo struct {
	EventType      string `json:"eventType"`      //值为1，表示是会话类型的消息
	ConvType       string `json:"convType"`       //会话具体类型：PERSON（二人会话数据）、TEAM（群聊数据）、	CUSTOM_PERSON（个人自定义系统通知）、CUSTOM_TEAM（群组自定义系统通知），字符串类型
	To             string `json:"to"`             //若convType为PERSON或CUSTOM_PERSON，则to为消息接收者的用户账号，字符串类型；若convType为TEAM或CUSTOM_TEAM，则to为tid，即群id，可转为Long型数据
	FromAccount    string `json:"fromAccount"`    //消息发送者的用户账号，字符串类型
	FromClientType string `json:"fromClientType"` //发送客户端类型： AOS、IOS、PC、WINPHONE、WEB、REST，字符串类型
	FromDeviceID   string `json:"fromDeviceId"`   //发送设备id，字符串类型
	FromNick       string `json:"fromNick"`       //发送方昵称，字符串类型
	MsgTimestamp   string `json:"msgTimestamp"`   //消息发送时间，字符串类型
	MsgType        string `json:"msgType"`        //会话具体类型PERSON、TEAM对应的通知消息类型:EXT、PICTURE、AUDIO、VIDEO、LOCATION 、NOTIFICATION、FILE、 //文件消息NETCALL_AUDIO、 //网络电话音频聊天 	NETCALL_VEDIO、 //网络电话视频聊天 	DATATUNNEL_NEW、 //新的数据通道请求通知 	TIPS、 //提示类型消息 	CUSTOM //自定义消息		会话具体类型CUSTOM_PERSON对应的通知消息类型：	FRIEND_ADD、 //加好友	FRIEND_DELETE、 //删除好友	CUSTOM_P2P_MSG、 //个人自定义系统通知		会话具体类型CUSTOM_TEAM对应的通知消息类型：	TEAM_APPLY、 //申请入群	TEAM_APPLY_REJECT、 //拒绝入群申请	TEAM_INVITE、 //邀请进群	TEAM_INVITE_REJECT、 //拒绝邀请	TLIST_UPDATE、 //群信息更新 	CUSTOM_TEAM_MSG、 //群组自定义系统通知
	Body           string `json:"body"`           //消息内容，字符串类型
	Attach         string `json:"attach"`         //附加消息，字符串类型
	MsgidClient    string `json:"msgidClient"`    //客户端生成的消息id，仅在convType为PERSON或TEAM含此字段，字符串类型
	MsgidServer    string `json:"msgidServer"`    //服务端生成的消息id，可转为Long型数据
	ResendFlag     string `json:"resendFlag"`     //重发标记：0不是重发, 1是重发。仅在convType为PERSON或TEAM时含此字段，可转为Integer类型数据
	CustomSafeFlag string `json:"customSafeFlag"` //自定义系统通知消息是否存离线:0：不存，1：存。仅在convType为CUSTOM_PERSON或CUSTOM_TEAM时含此字段，可转为Integer类型数据
	CustomApnsText string `json:"customApnsText"` //自定义系统通知消息推送文本。仅在convType为CUSTOM_PERSON或CUSTOM_TEAM时含此字段，字符串类型
	TMembers       string `json:"tMembers"`       //跟本次群操作有关的用户accid，仅在convType为TEAM或CUSTOM_TEAM时含此字段，字符串类型
	Ext            string `json:"ext"`            //消息扩展字段
	Antispam       string `json:"antispam"`       //标识是否被反垃圾，仅在被反垃圾时才有此字段，可转为Boolean类型数据
	YidunRes       string `json:"yidunRes"`       //易盾反垃圾的原始处理细节，只有接入了相关功能易盾反垃圾的应用才会有这个字段。
	Blacklist      string `json:"blacklist"`      //标识单聊消息是否黑名单，仅在消息发送方被拉黑时才有此字段，可转为 Boolean 类型数据
	Ip             string `json:"ip"`             //消息发送方的客户端 IP 地址(仅 SDK 发送的消息才有该字段)
}

type CopyEvent struct {
	EventType string `json:"eventType"` //为5，表示是实时音视频/白板时长类型事件
	Ext       string `json:"ext"`       //音视频发起时的自定义字段，可选，由用户指定

	*LogoutEventCopyInfo
	*MsgCopyInfo
}

func (c *ImClient) CheckMd5(req *http.Request) ([]byte, error) {
	if req == nil {
		return nil, errors.New("request 参数不能为空")
	}

	md5 := req.Header.Get("MD5")
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重要：重置 Body

	if Md5HashToHexString(bodyBytes) != md5 {
		return bodyBytes, fmt.Errorf("消息抄送内容被劫持,[md5]:%s,[body]:%s,[encodedBody]:%s", md5, string(bodyBytes), ShaHashToHexString(bodyBytes))
	}

	return bodyBytes, nil
}

func (c *ImClient) CheckSum(req *http.Request) error {
	if req == nil {
		return errors.New("request 参数不能为空")
	}

	checkSum := req.Header.Get("CheckSum")
	md5 := req.Header.Get("MD5")
	curTime := req.Header.Get("CurTime")

	recheck := ShaHashToHexStringFromString(c.AppSecret + md5 + curTime)
	if checkSum != recheck {
		return fmt.Errorf("CheckSum校验失败,[request-header-checkSum]:%s,[Checksum]:%s,[encodedChecksum]:%s", checkSum, c.AppSecret+md5+curTime, recheck)
	}
	return nil
}

func (c *ImClient) CheckSumMd5(req *http.Request) ([]byte, error) {
	if err := c.CheckSum(req); err != nil {
		return nil, err
	}

	return c.CheckMd5(req)
}

func (r *CopyEvent) UnmarshalJSON(data []byte) error {
	// 使用中间类型解除递归
	type Alias CopyEvent
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
	case EventTypeLogin, EventTypeLogout:
		if r.LogoutEventCopyInfo == nil {
			r.LogoutEventCopyInfo = &LogoutEventCopyInfo{}
		}
		return json.Unmarshal(data, r.LogoutEventCopyInfo)
	case EventTypeConversation:
		if r.MsgCopyInfo == nil {
			r.MsgCopyInfo = &MsgCopyInfo{}
		}
		err := json.Unmarshal(data, r.MsgCopyInfo)
		return err
	}

	return nil
}
