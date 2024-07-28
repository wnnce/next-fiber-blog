package service

import (
	"bytes"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/qiniu"
	"go-fiber-ent-web-layout/internal/tools/region"
	"go-fiber-ent-web-layout/internal/usercase"
	"go-fiber-ent-web-layout/pkg/optimize"
	"go-fiber-ent-web-layout/pkg/pool"
	"io"
	"log/slog"
	"math"
	"mime"
	"mime/multipart"
	"strings"
	"time"
)

type OtherService struct {
	repo usercase.IOtherRepo
}

func NewOtherService(repo usercase.IOtherRepo) usercase.IOtherService {
	return &OtherService{
		repo: repo,
	}
}

func (self *OtherService) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		slog.Error("文件打开失败", "err", err)
		return "", tools.FiberServerError("文件打开失败")
	}
	defer func() {
		_ = file.Close()
	}()
	originName := fileHeader.Filename
	var buffer *bytes.Buffer
	var suffix string
	if optimize.CheckWebpFormatSupport(originName) {
		if buffer, err = optimize.CompressImageWithWebp(file, false, 95); err == nil {
			suffix = ".webp"
		}
	}
	if buffer == nil {
		buffer = new(bytes.Buffer)
		_, _ = io.Copy(buffer, file)
		suffix = originName[strings.LastIndexByte(originName, '.'):]
	}
	sign := optimize.ComputeMd5(buffer.Bytes())
	// 查询图片是否已经上传
	uploadFile, err := self.repo.QueryFileByMd5(sign)
	if err != nil {
		slog.Error("检查图片是否上传失败", "err", err)
		return "", tools.FiberServerError("检查图片是否上传失败")
	}
	if uploadFile != nil {
		return uploadFile.FilePath, err
	}
	newFileName := sign + suffix
	fileSize := int64(buffer.Len())
	filePath := self.generateUploadPath("images/", newFileName)
	if err = qiniu.Upload(filePath, buffer); err != nil {
		slog.Error("七牛云文件上传失败", "err", err)
		return "", tools.FiberServerError("上传失败")
	}
	filePath = "/b-oss/" + filePath
	// 异步保存文件上传记录
	pool.Go(func() {
		fileType := mime.TypeByExtension(suffix)
		self.repo.SaveFileRecord(&usercase.UploadFile{
			FileMd5:    sign,
			OriginName: originName,
			FileName:   newFileName,
			FilePath:   filePath,
			FileSize:   fileSize,
			FileType:   fileType,
		})
	})
	buffer.Reset()
	return filePath, nil

}

func (self *OtherService) UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		slog.Error("文件打开失败", "err", err)
		return "", tools.FiberServerError("文件打开失败")
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("读取文件失败", "err", err)
		return "", tools.FiberServerError("文件读取失败")
	}
	sign := optimize.ComputeMd5(fileBytes)
	originName := fileHeader.Filename
	suffix := originName[strings.LastIndexByte(originName, '.'):]
	newFileName := sign + suffix
	filePath := self.generateUploadPath("files/", newFileName)
	fileSize := fileHeader.Size
	uploadFile, err := self.repo.QueryFileByMd5(sign)
	if uploadFile != nil && err == nil {
		return uploadFile.FilePath, nil
	}
	if err = qiniu.Upload(filePath, file); err != nil {
		slog.Error("七牛云文件上传失败", "err", err)
		return "", tools.FiberServerError("上传失败")
	}
	filePath = "/b-oss/" + filePath
	pool.Go(func() {
		fileType := mime.TypeByExtension(suffix)
		self.repo.SaveFileRecord(&usercase.UploadFile{
			FileMd5:    sign,
			OriginName: originName,
			FileName:   newFileName,
			FilePath:   filePath,
			FileSize:   fileSize,
			FileType:   fileType,
		})
	})
	return filePath, nil
}

func (self *OtherService) DeleteFile(filename string) {
	if err := self.repo.DeleteFileByName(filename); err != nil {
		slog.Error("删除文件上传记录失败", "err", err)
		return
	}
	pool.Go(func() {
		_ = qiniu.Remove(filename)
	})
}

func (self *OtherService) TraceLogin(record *usercase.LoginLog) {
	pool.Go(func() {
		location := region.SearchLocation(record.LoginIP)
		record.Location = location
		self.repo.SaveLoginRecord(record)
	})
}

func (self *OtherService) TraceAccess(referee, ip, ua string) {
	// 异步保存
	pool.Go(func() {
		location := region.SearchLocation(ip)
		self.repo.SaveAccessRecord(&usercase.AccessLog{
			Location: location,
			AccessIp: ip,
			AccessUa: ua,
			Referee:  referee,
		})
	})
}

func (self *OtherService) PageLogin(query *usercase.LoginLogQueryForm) (*usercase.PageData[usercase.LoginLog], error) {
	records, total, err := self.repo.PageLoginRecord(query)
	if err != nil {
		slog.Error("获取登录日志分页列表失败", "err", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	pages := int(math.Ceil(float64(total) / float64(query.Size)))
	return &usercase.PageData[usercase.LoginLog]{
		Current: query.Page,
		Size:    query.Size,
		Pages:   pages,
		Total:   total,
		Records: records,
	}, nil
}

func (self *OtherService) PageAccess(query *usercase.AccessLogQueryForm) (*usercase.PageData[usercase.AccessLog], error) {
	records, total, err := self.repo.PageAccessRecord(query)
	if err != nil {
		slog.Error("获取访问日志分页列表失败", "err", err.Error())
		return nil, tools.FiberServerError("查询失败")
	}
	pages := int(math.Ceil(float64(total) / float64(query.Size)))
	return &usercase.PageData[usercase.AccessLog]{
		Current: query.Page,
		Size:    query.Size,
		Pages:   pages,
		Total:   total,
		Records: records,
	}, nil
}

func (self *OtherService) generateUploadPath(prefix, fileName string) string {
	datePath := time.Now().Format("2006/0102/")
	return prefix + datePath + fileName
}
