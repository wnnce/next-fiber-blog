package service

import (
	"bytes"
	"context"
	"go-fiber-ent-web-layout/internal/data"
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
	"sync"
	"time"
)

const SiteConfigurationCacheKey = "BLOG:site_configuration"

type OtherService struct {
	repo          usercase.IOtherRepo
	redisTemplate *data.RedisTemplate
}

func NewOtherService(repo usercase.IOtherRepo, redisTemplate *data.RedisTemplate) usercase.IOtherService {
	return &OtherService{
		repo:          repo,
		redisTemplate: redisTemplate,
	}
}

func (self *OtherService) SiteConfiguration() map[string]usercase.SiteConfigurationItem {
	config, err := data.RedisGetStruct[map[string]usercase.SiteConfigurationItem](context.Background(), SiteConfigurationCacheKey)
	if err == nil && config != nil {
		return config
	}
	return usercase.DefaultSiteConfiguration
}

func (self *OtherService) UpdateSiteConfiguration(config map[string]usercase.SiteConfigurationItem) error {
	err := self.redisTemplate.Set(context.Background(), SiteConfigurationCacheKey, config, time.Duration(math.MaxInt64))
	if err != nil {
		slog.Error("更新站点配置失败", "error", err.Error())
		return tools.FiberServerError("更新失败")
	}
	return nil
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

func (self *OtherService) SiteStats() (usercase.SiteStats, error) {
	stats, err := self.repo.SiteStats()
	if err != nil {
		slog.Error("查询站点统计数据失败", "error", err.Error())
		return stats, tools.FiberServerError("查询失败")
	}
	return stats, nil
}

func (self *OtherService) AdminIndexStats() (*usercase.AdminIndexStats, error) {
	var resultMap sync.Map
	var wg sync.WaitGroup
	wg.Add(5)
	// 异步同时查询所有数据
	pool.Go(func() {
		defer wg.Done()
		stats, err := self.repo.AdminIndexStats()
		if err != nil {
			slog.Error("查询后台首页统计数据失败", "error", err.Error())
			return
		}
		resultMap.Store("stats", stats)
	})
	pool.Go(func() {
		defer wg.Done()
		accessArray, err := self.repo.AccessStatsArray()
		if err != nil {
			slog.Error("查询访问记录数据失败", "error", err.Error())
			return
		}
		resultMap.Store("access", accessArray)
	})
	pool.Go(func() {
		defer wg.Done()
		commentArray, err := self.repo.CommentStatsArray()
		if err != nil {
			slog.Error("查询评论统计数据失败", "error", err.Error())
			return
		}
		resultMap.Store("comment", commentArray)
	})
	pool.Go(func() {
		defer wg.Done()
		userArray, err := self.repo.UserStatsArray()
		if err != nil {
			slog.Error("查询用户统计数据失败", "error", err.Error())
			return
		}
		resultMap.Store("user", userArray)
	})
	pool.Go(func() {
		defer wg.Done()
		articleArray, err := self.repo.ArticleStatsArray()
		if err != nil {
			slog.Error("查询文章统计数据失败", "error", err.Error())
			return
		}
		resultMap.Store("article", articleArray)
	})
	// 等待所有协程执行完成
	wg.Wait()
	// 拼装数据
	stats, ok := resultMap.Load("stats")
	if !ok {
		return nil, tools.FiberServerError("查询统计数据失败")
	}
	result := stats.(usercase.AdminIndexStats)
	if accessArray, ok := resultMap.Load("access"); ok {
		result.AccessArray = accessArray.([]usercase.DayStats)
	}
	if commentArray, ok := resultMap.Load("comment"); ok {
		result.CommentArray = commentArray.([]usercase.DayStats)
	}
	if userArray, ok := resultMap.Load("user"); ok {
		result.UserArray = userArray.([]usercase.DayStats)
	}
	if articleArray, ok := resultMap.Load("article"); ok {
		result.ArticleArray = articleArray.([]usercase.DayStats)
	}
	return &result, nil
}

func (self *OtherService) generateUploadPath(prefix, fileName string) string {
	datePath := time.Now().Format("2006/0102/")
	return prefix + datePath + fileName
}
