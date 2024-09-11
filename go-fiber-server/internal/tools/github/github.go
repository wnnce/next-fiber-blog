package github

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
	"go-fiber-ent-web-layout/internal/conf"
)

// Config Github配置
type Config struct {
	ClientId     string `json:"clientId" yaml:"client-id"`
	ClientSecret string `json:"clientSecret" yaml:"client-secret"`
	Proxy        string `json:"proxy" yaml:"proxy"`
}

// Email Github邮箱详细信息
type Email struct {
	Email      string `json:"email"`                // 邮箱地址
	Primary    bool   `json:"primary"`              // 是否主要邮箱
	Verified   bool   `json:"verified"`             // 邮箱是否验证
	Visibility string `json:"visibility,omitempty"` // 邮箱是否可见
}

const (
	tokenUrl    = "https://github.com/login/oauth/access_token"
	userInfoUrl = "https://api.github.com/user"
	emailsUrl   = "https://api.github.com/user/emails"
)

var (
	clientId     string
	clientSecret string
	proxy        string
	client       *fasthttp.Client
)

func init() {
	conf.CollectConfigReader(func(config *conf.Bootstrap) {
		clientId = config.Github.ClientId
		clientSecret = config.Github.ClientSecret
		proxy = config.Github.Proxy
	})
	client = &fasthttp.Client{
		ReadBufferSize: 8192,
	}
}

func AccessToken(code string) (string, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	request.Header.Add(fiber.HeaderAccept, "application/json")
	request.Header.SetMethod(fiber.MethodPost)
	request.SetRequestURI(fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s", proxy+tokenUrl, clientId, clientSecret, code))
	if err := client.Do(request, response); err != nil {
		return "", err
	}
	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()
	result := make(map[string]string)
	if err := sonic.Unmarshal(response.Body(), &result); err != nil {
		return "", err
	}
	accessToken, ok := result["access_token"]
	if !ok {
		return "", fmt.Errorf("not AccessToken")
	}
	return accessToken, nil
}

func UserProfile(accessToken string) (map[string]string, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	request.Header.SetMethod(fiber.MethodGet)
	request.Header.Add(fiber.HeaderAccept, "application/vnd.github+json")
	request.Header.Add(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", accessToken))
	request.SetRequestURI(proxy + userInfoUrl)
	if err := client.Do(request, response); err != nil {
		return nil, err
	}
	defer func() {
		fasthttp.ReleaseResponse(response)
		fasthttp.ReleaseRequest(request)
	}()
	if fiber.StatusOK != response.StatusCode() {
		return nil, fmt.Errorf("request field response status: %d", response.StatusCode())
	}
	result := make(map[string]string)
	if err := sonic.Unmarshal(response.Body(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func UserEmails(accessToken string) ([]*Email, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	request.Header.SetMethod(fiber.MethodGet)
	request.Header.Add(fiber.HeaderAccept, "application/vnd.github+json")
	request.Header.Add(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", accessToken))
	request.SetRequestURI(proxy + emailsUrl)
	if err := client.Do(request, response); err != nil {
		return nil, err
	}
	defer func() {
		fasthttp.ReleaseResponse(response)
		fasthttp.ReleaseRequest(request)
	}()
	if fiber.StatusOK != response.StatusCode() {
		return nil, fmt.Errorf("request field response status: %d", response.StatusCode())
	}
	emails := make([]*Email, 0)
	if err := sonic.Unmarshal(response.Body(), &emails); err != nil {
		return nil, err
	}
	return emails, nil
}
