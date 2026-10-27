package yunxin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"log/slog"

	"github.com/google/go-querystring/query"
)

const neteaseBaseURL = "https://open.yunxinapi.com/im/v2"

// defaultTimeout 默认 HTTP 请求超时时间，可通过 WithHttpClient 覆盖
const defaultTimeout = 30 * time.Second

// Log 日志抽象接口。注意：*slog.Logger 本身的 Debug 方法不接收 context，
// 需通过 slogLog 适配器将其 DebugContext 适配为 Log.Debug 签名。
type Log interface {
	Debug(ctx context.Context, msg string, args ...any)
}

// slogLog 将 *slog.Logger 适配为 Log 接口：
// Log.Debug(ctx, msg, args...)  ->  slog.Logger.DebugContext(ctx, msg, args...)
type slogLog struct {
	*slog.Logger
}

func (s slogLog) Debug(ctx context.Context, msg string, args ...any) {
	s.Logger.DebugContext(ctx, msg, args...)
}

type ResponseBase struct {
	Code int    `json:"code"` // 状态码，200 表示成功
	Msg  string `json:"msg"`  // 提示信息：失败时返回错误信息，成功为 "success"
}

type Response[T any] struct {
	ResponseBase
	Data *T `json:"data,omitempty"` // 返回数据对象，请求失败则为空对象
}

type Failed struct {
	AccountId string `json:"account_id"` //查询失败的 IM 账号 ID。
	ErrorCode int    `json:"error_code"` //	查询失败的错误码。
	ErrorMsg  string `json:"error_msg"`  //查询失败的提示信息。
}

// ImClient .
type ImClient struct {
	AppKey    string
	AppSecret string

	DryRun bool // 是否开启 Dry Run：开启时不真正发送 HTTP 请求

	log        Log
	httpClient *http.Client
	RetryMax   int           // 最大重试次数
	RetryWait  time.Duration // 初始等待时间（用于指数退避）
}

type Option func(c *ImClient)

func WithLogger(logger Log) Option {
	return func(c *ImClient) {
		c.log = logger
	}
}

func WithHttpClient(client *http.Client) Option {
	return func(c *ImClient) {
		c.httpClient = client
	}
}

// WithRetryMax 设置最大重试次数。
// 注意：重试仅针对「网络错误 / HTTP 非 200 状态码」等传输层失败；
// 业务失败（响应 code != 200）不会触发重试，将直接返回错误。
func WithRetryMax(max int) Option {
	return func(c *ImClient) {
		c.RetryMax = max
	}
}

// WithRetryWait 设置重试的初始等待时间（用于指数退避）。
func WithRetryWait(wait time.Duration) Option {
	return func(c *ImClient) {
		c.RetryWait = wait
	}
}

func WithDryRun(enabled bool) Option {
	return func(c *ImClient) {
		c.DryRun = enabled
	}
}

// CreateImClient  创建im客户端
func CreateImClient(appKey, appSecret string, opts ...Option) *ImClient {
	c := &ImClient{
		AppKey:    appKey,
		AppSecret: appSecret,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: defaultTimeout}
	}

	if c.log == nil {
		c.log = slogLog{slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,            // 记录日志位置
			Level:     slog.LevelDebug, // 设置日志级别
		}))}
	}

	return c
}

func (c *ImClient) buildRequest(ctx context.Context, method, urlAddr string, param any) (*http.Request, error) {
	var payloadBytes []byte
	if method != "GET" && param != nil {
		var err error
		payloadBytes, err = json.Marshal(param)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	} else if method == "GET" && param != nil {
		values, err := query.Values(param)
		if err != nil {
			return nil, fmt.Errorf("failed to encode query parameters: %w", err)
		}
		urlAddr = urlAddr + "?" + values.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, urlAddr, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 设置请求头（Nonce 每次请求独立生成，配合 CurTime + CheckSum 保证鉴权安全）
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json;charset=utf-8")
	headers.Set("Accept", "application/json;charset=utf-8")
	headers.Set("AppKey", c.AppKey)
	nonce := RandStringBytesMaskImprSrc(64)
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	headers.Set("Nonce", nonce)
	headers.Set("CurTime", timeStamp)
	headers.Set("CheckSum", ShaHashToHexStringFromString(c.AppSecret+nonce+timeStamp))
	req.Header = headers

	// 注意：Debug 级别会打印完整请求头与请求体（可能包含业务消息内容），
	// 生产环境请勿开启 Debug 日志，或在自定义 Logger 中对 header/body 做脱敏处理。
	// 采用惰性参数传递，仅当 Debug 级别启用时才格式化，避免热路径上的无谓开销。
	c.log.Debug(ctx, "curl",
		"method", req.Method,
		"url", req.URL.String(),
		"header", headers,
		"body", string(payloadBytes),
	)

	return req, nil
}

// sendRequestWithRetry 执行带指数退避重试的 HTTP 请求。
// 重试边界：仅当请求因网络错误失败，或返回非 200 的 HTTP 状态码时触发重试；
// 业务层失败（code != 200）由调用方在 JSON 反序列化后判定，不在此处重试。
// 重试复用同一 *http.Request，依赖 req.GetBody 在每次重试前重置请求体。
func (c *ImClient) sendRequestWithRetry(ctx context.Context, req *http.Request) ([]byte, *http.Response, error) {
	var resp *http.Response
	var respBody []byte
	var err error

	maxRetries := c.RetryMax
	waitTime := c.RetryWait

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			if waitTime == 0 {
				waitTime = time.Second // 默认初始等待时间为 1 秒
			}
			time.Sleep(waitTime)
			waitTime *= 2 // 指数退避

			c.log.Debug(ctx, "Retrying request", "attempt", attempt+1, "wait", waitTime)
		}

		// 重置请求体，避免重试时因 body 已被消费而发送空请求
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to reset request body before retry: %w", err)
			}
			req.Body = body
		}

		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		// 关闭非成功响应的 body，避免重试过程中连接泄漏
		if resp != nil {
			resp.Body.Close()
		}

		if attempt < maxRetries {
			continue
		}

		if err == nil {
			err = fmt.Errorf("HTTP status code not OK: %d", resp.StatusCode)
		}
		return nil, nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, err)
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, resp, nil
}

func (c *ImClient) Curl(ctx context.Context, method, urlAddr string, param any) (*http.Response, []byte, error) {
	req, err := c.buildRequest(ctx, method, urlAddr, param)
	if err != nil {
		return nil, nil, err
	}

	if c.DryRun {
		c.log.Debug(ctx, "dry run: skip request", "url", req.URL.String())
		return nil, nil, fmt.Errorf("dry run: no real request sent to %s", req.URL.String())
	}

	respBody, resp, err := c.sendRequestWithRetry(ctx, req)
	if err != nil {
		return resp, respBody, err
	}

	return resp, respBody, nil
}

func (res *Response[T]) Curl(ctx context.Context, c *ImClient, method, urlAddr string, param any) (*http.Response, []byte, error) {
	resp, respBody, err := c.Curl(ctx, method, urlAddr, param)

	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return resp, respBody, err
	}

	if res.Code != 200 {
		err = fmt.Errorf("code:%v,msg:%s", res.Code, res.Msg)
		return resp, respBody, err
	}

	return resp, respBody, nil
}