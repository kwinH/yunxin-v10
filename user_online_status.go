package yunxin

import (
	"context"
	"fmt"
)

type OnlineStatusSuccess struct {
	AccountId    string `json:"account_id"` //IM 账号 ID。
	OnlineStatus []struct {
		ClientType int   `json:"client_type"` //登录的客户端类型。 1：AOS。2：IOS：4：PC。16：WEB。64：MAC。65：HARMONY。
		LoginTime  int64 `json:"login_time"`  //登录时间。
	} `json:"online_status"` // 用户在线状态。
}

type OnlineStatusResponse struct {
	SuccessList []OnlineStatusSuccess `json:"success_list"` //查询成功的用户在线状态信息列表。
	FailedList  []Failed              `json:"failed_list"`  //查询失败的 IM 账号列表。
}

func (c *ImClient) GetUserOnlineStatus(ctx context.Context, accountIds ...string) (res Response[OnlineStatusResponse], err error) {
	param := map[string]any{
		"account_ids": accountIds,
	}
	_, _, err = res.Curl(ctx, c, "POST", fmt.Sprintf("%s/users/actions/online_status", neteaseBaseURL), param)

	return
}
