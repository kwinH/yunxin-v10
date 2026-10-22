package yunxin

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"strings"
)

// RandStringBytesMaskImprSrc 使用 crypto/rand 生成指定长度的随机字符串。
func RandStringBytesMaskImprSrc(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if n <= 0 {
		return ""
	}

	rb := make([]byte, n)
	if _, err := rand.Read(rb); err != nil {
		// crypto/rand 失败的概率极低，失败时无法安全降级，直接 panic
		panic(err)
	}

	out := make([]byte, n)
	for i := range rb {
		out[i] = alphabet[int(rb[i])%len(alphabet)]
	}
	return string(out)
}

// ShaHashToHexString SHA1加密字符串，并将加密结果转成16进制字符串
func ShaHashToHexString(bv []byte) string {
	hasher := sha1.New()
	hasher.Write(bv)
	return strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))
}

// ShaHashToHexStringFromString .
func ShaHashToHexStringFromString(src string) string {
	return ShaHashToHexString([]byte(src))
}

// Md5HashToHexString .
func Md5HashToHexString(bv []byte) string {
	hasher := md5.New()
	hasher.Write(bv)
	return strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))
}

func MarshalToString(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	return string(data), err
}

// Bool 返回 bool 值的指针。
// 用于向带 omitempty 的 *bool 配置字段显式传递 true/false，
// 避免零值 false 被 JSON 序列化时丢弃。nil 表示「不设置」该字段。
func Bool(v bool) *bool { return &v }
