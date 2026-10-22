# yunxin

yunxin-v10 是用 GO 语言实现的网易云信服务端 V10 系列 API 封装，目前实现了常用的大部分功能。如有其他需要或功能失效，欢迎提 issue 反馈，也欢迎提 PR 完善代码。

![](https://img.shields.io/badge/language-golang-blue.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 安装

```
go get -u github.com/kwinh/yunxin-v10
```

## 导入

```go
import "github.com/kwinh/yunxin-v10"
```

## 使用

### 创建客户端

```go
ctx := context.Background()
client := yunxin.CreateImClient("AppKey", "AppSecret")
```

客户端支持通过函数式 Option 注入自定义配置：

```go
client := yunxin.CreateImClient(
    "AppKey",
    "AppSecret",
    yunxin.WithLogger(customLogger),                     // 自定义日志
    yunxin.WithHttpClient(customHTTPClient),             // 自定义 HTTP Client
    yunxin.WithRetryMax(3),                               // 最大重试次数（指数退避）
    yunxin.WithRetryWait(time.Second),                    // 重试初始等待时间
    yunxin.WithDryRun(true),                              // 开启 Dry Run：不真正发送请求
)
```

### 创建 IM 账号

`CreateAccount` 在 `CreateIMAccountRequest` 中指定 `AccountId`（必填）以及可选的 `Token`、`Configuration`（账号配置）、`UserInfo`（用户信息）。

```go
res, err := client.CreateAccount(ctx, yunxin.CreateIMAccountRequest{
    AccountId: "acc-001",
    Configuration: yunxin.AccountConfiguration{
        Enabled: yunxin.Bool(true), // *bool：nil=不设置，显式传 true/false 均可
    },
    UserInfo:  yunxin.UserInfo{Name: "张三", Gender: 1},
})
if err != nil {
    log.Printf("创建账号失败: %v", err)
}
// res.Data.AccountId ... (res.Data) 请求失败时忽略
```

> **布尔配置字段说明**：部分标注为「默认 true」的开关字段（如 `AccountConfiguration.Enabled`、`AccountConfiguration.PushEnabledWhenDesktopOnline`、`TargetOption.Inclusive`、`NotificationConfig.*` 等）类型为 `*bool`，用于显式区分「不设置」与「传 `false`」两种语义：
> - **不设置**（`nil`）：字段不发送，采用云信默认值；
> - **`yunxin.Bool(true)` / `yunxin.Bool(false)`**：显式发送对应布尔值，避免 `false` 被省略。
>
> 其余未标注「默认 true」的开关字段保持普通 `bool` 类型，直接赋值 `true`/`false` 即可。

### 发送消息

`SendMessage` 通过 `ownerId`、`conversationType`、`to` 拼接会话 ID：
- `conversationType`：1=单聊，2=高级群，3=超大群
- `Message.MessageType`：0=文本，1=图片，2=语音，3=视频，4=地理位置，6=文件，100=自定义
- 携带可选的消息配置（`MessageConfig`）、推送、抄送、反垃圾等

```go
res, err := client.SendMessage(ctx, "1", 1, "3", yunxin.SendMessageRequest{
    Message: yunxin.Message{
        MessageType: 0,
        Text:        "hello world",
    },
})
```

自定义消息（`message_type=100`，`attachment` 为自定义内容）：

```go
res, err := client.SendMessage(ctx, "1", 1, "3", yunxin.SendMessageRequest{
    Message: yunxin.Message{
        MessageType: 100,
        SubType:     1,
        Attachment:  map[string]any{"action": "open_doc"},
    },
})
```

### 发送流式消息

```go
res, err := client.SendStreamMessage(ctx, "1", 1, "3", yunxin.StreamMessageRequest{
    Message: yunxin.StreamMessage{
        MessageType: 0,
        Attachment:  yunxin.Attachment{Text: "分片内容", Index: 0, Finish: 1},
    },
})
```

### 发送自定义系统通知

```go
res, err := client.SendCustomNotification(ctx, yunxin.CustomNotificationRequest{
    Type:       1, // 1=单聊，2=高级群，3=超大群
    SenderId:   "1",
    ReceiverId: "3",
    Content:    `{"code":1,"msg":"通知内容"}`,
    NotificationConfig: &yunxin.NotificationConfig{
        OfflineEnabled: yunxin.Bool(false), // 显式传 false：不存离线消息
    },
})
```

### 查询历史消息

```go
res, err := client.GetMessagesPage(ctx, "1", 1, "3", yunxin.QueryMessagesRequest{
    BeginTime:  now.UnixMilli() - 3600*1000,
    EndTime:    now.UnixMilli(),
    Limit:      20,
})
```

### 修改消息

```go
res, err := client.ModifyMsg(ctx, yunxin.ModifyMsgRequest{
    Operator: "1",
    Type:     1, // 1=单聊，2=高级群，3=超大群
    Message: yunxin.UpdateMessage{
        MessageServerId: serverId,
        ReceiverId:      "3",
        Time:            sendTime,
        SenderId:        "1",
        MessageType:     0,
        Text:            "修改后的内容",
    },
})
```

### 账号 / 用户管理

```go
// 查询账号信息
res, err := client.AccountInfo(ctx, "1")

// 封禁 / 解封账号
res, err := client.AccountDisable(ctx, "1", yunxin.DisableIMAccountRequest{Enabled: false, NeedKick: true})

// 刷新 Token
res, err := client.RefreshToken(ctx, "1")

// 查询单个用户信息
res, err := client.GetUserInfo(ctx, "1")

// 批量查询用户信息
res, err := client.GetUserInfos(ctx, "1", "2", "3")

// 更新用户信息
res, err := client.UpdateUserInfo(ctx, "1", yunxin.UpdateUserInfoRequest{Name: "新昵称", Avatar: "https://..."})

// 查询用户在线状态
res, err := client.GetUserOnlineStatus(ctx, "1", "2")
```

## 已实现功能

- [x] 账号：创建账号、查询账号信息、封禁/解封、刷新 Token
- [x] 用户：查询（单个/批量）、更新用户信息、查询在线状态
- [x] 消息：发送文本/图片/语音/视频/自定义消息、批量消息配置、流式消息
- [x] 消息管理：分页查询单聊历史、修改（撤回/编辑）消息
- [x] 系统通知：发送自定义系统通知、系统系统通知
- [x] 基础设施：请求鉴权（AppKey/Nonce/CurTime/CheckSum）、指数退避重试、Dry Run、自定义日志

## License

yunxin 使用 [MIT](https://opensource.org/licenses/MIT) 开源协议