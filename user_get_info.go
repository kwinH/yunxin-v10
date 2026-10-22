package yunxin

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type UserInfoResponse struct {
	AccountId string `json:"account_id,omitempty"` // IM 账号 ID
	UserInfo
}

type UserInfosResponse struct {
	SuccessList []UserInfoResponse `json:"success_list"` //查询成功的用户信息列表。
	FailedList  []Failed           `json:"failed_list"`  //查询失败的 IM 账号列表。
}

func (c *ImClient) GetUserInfo(ctx context.Context, accountId string) (res Response[UserInfoResponse], err error) {
	_, _, err = res.Curl(ctx, c, "GET", fmt.Sprintf("%s/users/%s", neteaseBaseURL, accountId), nil)

	return
}

func (c *ImClient) GetUserInfos(ctx context.Context, accountIds ...string) (res Response[UserInfosResponse], err error) {
	encoded := make([]string, len(accountIds))
	for i, id := range accountIds {
		encoded[i] = url.QueryEscape(id)
	}
	addr := fmt.Sprintf("%s/users?account_ids=%s", neteaseBaseURL, strings.Join(encoded, ","))

	_, _, err = res.Curl(ctx, c, "GET", addr, nil)

	return
}
