package service

import (
	"bytes"
	"go-fiber-ent-web-layout/internal/tools"
	"go-fiber-ent-web-layout/internal/tools/pool"
	"go-fiber-ent-web-layout/internal/tools/qiniu"
	"go-fiber-ent-web-layout/internal/tools/region"
	"go-fiber-ent-web-layout/internal/usercase"
	"io"
	"log/slog"
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

func (ot *OtherService) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
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
	if tools.CheckWebpFormatSupport(originName) {
		if buffer, err = tools.ImageToWebp(file, false, 95); err == nil {
			suffix = ".webp"
		}
	}
	if buffer == nil {
		buffer = new(bytes.Buffer)
		_, _ = io.Copy(buffer, file)
		suffix = originName[strings.LastIndexByte(originName, '.'):]
	}
	sign := tools.ComputeMd5(buffer.Bytes())
	// 查询图片是否已经上传
	uploadFile, err := ot.repo.QueryFileByMd5(sign)
	if err != nil {
		slog.Error("检查图片是否上传失败", "err", err)
		return "", tools.FiberServerError("检查图片是否上传失败")
	}
	if uploadFile != nil {
		return uploadFile.FilePath, err
	}
	newFileName := sign + suffix
	fileSize := int64(buffer.Len())
	filePath := ot.generateUploadPath("images/", newFileName)
	if err = qiniu.Upload(filePath, buffer); err != nil {
		slog.Error("七牛云文件上传失败", "err", err)
		return "", tools.FiberServerError("上传失败")
	}
	// 异步保存文件上传记录
	pool.Go(func() {
		fileType := mime.TypeByExtension(suffix)
		ot.repo.SaveFileRecord(&usercase.UploadFile{
			FileMd5:    sign,
			OriginName: originName,
			FileName:   newFileName,
			FilePath:   "b-oss/" + filePath,
			FileSize:   fileSize,
			FileType:   fileType,
		})
	})
	buffer.Reset()
	return filePath, nil

}

func (ot *OtherService) UploadFile(fileHeader *multipart.FileHeader) (string, error) {
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
	sign := tools.ComputeMd5(fileBytes)
	originName := fileHeader.Filename
	suffix := originName[strings.LastIndexByte(originName, '.'):]
	newFileName := sign + suffix
	filePath := ot.generateUploadPath("files/", newFileName)
	fileSize := fileHeader.Size
	uploadFile, err := ot.repo.QueryFileByMd5(sign)
	if uploadFile != nil && err == nil {
		return uploadFile.FilePath, nil
	}
	if err = qiniu.Upload(filePath, file); err != nil {
		slog.Error("七牛云文件上传失败", "err", err)
		return "", tools.FiberServerError("上传失败")
	}
	pool.Go(func() {
		fileType := mime.TypeByExtension(suffix)
		ot.repo.SaveFileRecord(&usercase.UploadFile{
			FileMd5:    sign,
			OriginName: originName,
			FileName:   newFileName,
			FilePath:   "b-oss/" + filePath,
			FileSize:   fileSize,
			FileType:   fileType,
		})
	})
	return filePath, nil
}

func (ot *OtherService) DeleteFile(filename string) {
	if err := ot.repo.DeleteFileByName(filename); err != nil {
		slog.Error("删除文件上传记录失败", "err", err)
		return
	}
	pool.Go(func() {
		_ = qiniu.Remove(filename)
	})
}

func (ots *OtherService) TraceLogin(record *usercase.LoginLog) {
	pool.Go(func() {
		location := region.SearchLocation(record.LoginIP)
		record.Location = location
		ots.repo.SaveLoginRecord(record)
	})
}

func (ot *OtherService) TraceAccess(referee, ip, ua string) {
	// 异步保存
	pool.Go(func() {
		location := region.SearchLocation(ip)
		ot.repo.SaveAccessRecord(&usercase.AccessLog{
			Location: location,
			AccessIp: ip,
			AccessUa: ua,
			Referee:  referee,
		})
	})
}

func (ot *OtherService) generateUploadPath(prefix, fileName string) string {
	datePath := time.Now().Format("2006/0102/")
	return prefix + datePath + fileName
}
