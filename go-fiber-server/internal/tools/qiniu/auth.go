package qiniu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/bytedance/sonic"
	"net/url"
	"sort"
	"strings"
	"time"
)

const uploadReturnBody = `{"key": $(key), "size": $(size),"hash": $(etag), "w": $(imageInfo.width), "h": $(imageInfo.height)}`

type policy struct {
	Scope      string `json:"scope"`
	Deadline   int64  `json:"deadline"`
	ReturnBody string `json:"returnBody"`
}

// 生成上传需要的Token
// fileName 文件名称
// expire Token的过期时间
func makeUploadToken(fileName string, expire time.Duration) string {
	scope := qiniuConfig.Bucket + ":" + fileName
	expireTime := (time.Now().UnixMilli() + expire.Milliseconds()) / 1000
	body := &policy{
		Scope:      scope,
		Deadline:   expireTime,
		ReturnBody: uploadReturnBody,
	}
	bodyBytes, _ := sonic.Marshal(body)
	base64EncodeBody := base64.URLEncoding.EncodeToString(bodyBytes)
	sign := hmacSha1Signature(qiniuConfig.SecretKey, base64EncodeBody)
	base64EncodeSign := base64.URLEncoding.EncodeToString(sign)
	return qiniuConfig.AccessKey + ":" + base64EncodeSign + ":" + base64EncodeBody
}

// 生成管理文件所需要的管理凭证
// method 请求方式
// uri 请求的路径
// headers 自定义请求头
// qiniuHeaders 七牛云需要的请求头
func makeManageToken(method string, uri *url.URL, headers, qiniuHeaders map[string]string) string {
	var builder strings.Builder
	builder.WriteString(strings.ToUpper(method))
	builder.WriteByte(' ')
	builder.WriteString(uri.Path)
	builder.WriteString("\nHost: ")
	builder.WriteString(uri.Host)
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			builder.WriteByte('\n')
			builder.WriteString(k + ": " + v)
		}
	}
	if qiniuHeaders != nil && len(qiniuHeaders) > 0 {
		keys := make([]string, 0, len(qiniuHeaders))
		for k, _ := range qiniuHeaders {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value := qiniuHeaders[key]
			builder.WriteByte('\n')
			builder.WriteString(key + ": " + value)
		}
	}
	builder.WriteString("\n\n")
	sign := hmacSha1Signature(qiniuConfig.SecretKey, builder.String())
	token := base64.URLEncoding.EncodeToString(sign)
	return "Qiniu " + qiniuConfig.AccessKey + ":" + token
}

// hmacSha1 签名
func hmacSha1Signature(key, value string) []byte {
	hasher := hmac.New(sha1.New, []byte(key))
	hasher.Write([]byte(value))
	return hasher.Sum(nil)
}
