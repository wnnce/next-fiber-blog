package manage

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"go-fiber-ent-web-layout/internal/tools/clog"
	"go-fiber-ent-web-layout/internal/tools/res"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// LoggerPush 日志SSE推送
func LoggerPush(c fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// 推送日志的间隔时间 降低网络io操作
		interval, _ := strconv.ParseInt(c.Params("interval", "500"), 10, 0)
		_, _ = w.Write(make([]byte, 0))
		if err := w.Flush(); err != nil {
			return
		}
		ch := make(chan []byte, 10)
		// 日志缓冲区
		var buff bytes.Buffer
		var mu sync.Mutex

		registerKey := clog.RegisterChan(ch)
		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
		defer func() {
			ticker.Stop()
			close(ch)
			clog.RemoveChan(registerKey)
		}()
	Loop:
		for {
			select {
			case <-ticker.C:
				if buff.Len() > 0 {
					mu.Lock()
					data := buff.String()
					buff.Reset()
					mu.Unlock()
					if _, err := fmt.Fprintf(w, "data: [%s]\n\n", data); err != nil {
						break Loop
					}
					if err := w.Flush(); err != nil {
						break Loop
					}
				}
			case value := <-ch:
				mu.Lock()
				// 如果不是第一条数据 那么添加逗号
				if buff.Len() > 0 {
					buff.WriteByte(',')
				}
				// 去掉换行符
				buff.Write(value[:len(value)-1])
				mu.Unlock()
			}
		}
	})
	return nil
}

type ApplicationMonitor struct {
	runtime.MemStats
	Hostname          string  `json:"hostname"`
	Platform          string  `json:"platform"`
	PlatformVersion   string  `json:"platformVersion"`
	CpuNumber         int     `json:"cpuNumber"`
	CpuPercent        float64 `json:"cpuPercent"`
	MemoryTotal       uint64  `json:"memoryTotal"`
	MemoryUsed        uint64  `json:"memoryUsed"`
	MemoryAvailable   uint64  `json:"memoryAvailable"`
	MemoryUsedPercent float64 `json:"memoryUsedPercent"`
}

// Monitor 指标监控
func Monitor(ctx fiber.Ctx) error {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)

	monitor := ApplicationMonitor{
		MemStats: *stats,
	}
	// 获取系统信息
	if hostInfo, err := host.Info(); err == nil {
		monitor.Hostname = hostInfo.Hostname
		monitor.Platform = hostInfo.Platform
		monitor.PlatformVersion = hostInfo.PlatformVersion
	}
	// 获取内存信息
	if memInfo, err := mem.VirtualMemory(); err == nil {
		monitor.MemoryTotal = memInfo.Total
		monitor.MemoryUsed = memInfo.Used
		monitor.MemoryAvailable = memInfo.Available
		monitor.MemoryUsedPercent = memInfo.UsedPercent
	}
	// 获取cpu核数
	if cpuCount, err := cpu.Counts(true); err == nil {
		monitor.CpuNumber = cpuCount
	}
	// 获取cpu使用率
	if cpuPercent, err := cpu.Percent(time.Second, true); err == nil {
		monitor.CpuPercent = cpuPercent[0]
	}
	return ctx.JSON(res.OkByData(&monitor))
}
