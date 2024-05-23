package qiniu

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v3/client"
	"io"
	"log/slog"
	"mime/multipart"
	"net/url"
	"strings"
	"time"
)

// Upload 上传文件
// fileName 文件上传的路径和名称
// reader 文件内容
func Upload(fileName string, reader io.Reader) error {
	uploadToken := makeUploadToken(fileName, 1*time.Hour)
	request := client.AcquireRequest()
	requestBody := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(requestBody)
	fileWriter, _ := multipartWriter.CreateFormFile("file", fileName)
	_, _ = io.Copy(fileWriter, reader)
	_ = multipartWriter.WriteField("key", fileName)
	_ = multipartWriter.WriteField("token", uploadToken)
	_ = multipartWriter.Close()
	request.SetRawBody(requestBody.Bytes())
	request.AddHeader("Content-Type", multipartWriter.FormDataContentType())
	response, err := request.Post(qiniuConfig.Region)
	if err != nil {
		return err
	}
	if response.StatusCode() != 200 || len(response.Body()) == 0 {
		return fmt.Errorf("request error status: %d", response.StatusCode())
	}
	return nil
}

// Remove 删除上传的文件
// fileName 文件上传的路径和名称
func Remove(fileName string) error {
	encodeName := base64.URLEncoding.EncodeToString([]byte(qiniuConfig.Bucket + ":" + fileName))
	uri, err := url.Parse("https://rs-z0.qiniuapi.com/delete/" + encodeName)
	if err != nil {
		slog.Error("删除文件失败，创建链接地址错误", "err", err.Error())
		return err
	}
	now := time.Now().UTC().Format("20060102T150105Z")
	headers := make(map[string]string, 2)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	qiniuHeaders := map[string]string{
		"X-Qiniu-Date": now,
	}
	token := makeManageToken("POST", uri, headers, qiniuHeaders)
	headers["Authorization"] = token
	request := client.AcquireRequest()
	response, err := request.SetHeaders(headers).AddHeader("X-Qiniu-Date", now).Post(uri.String())
	if err != nil {
		slog.Error("七牛云发送删除请求失败", "err", err)
		return err
	}
	str := string(response.Body())
	if !strings.Contains(str, "200") {
		slog.Error("七牛云文件删除失败", "response", str)
		return fmt.Errorf("删除文件响应失败")
	}
	return nil
}
