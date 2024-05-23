package manage

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go-fiber-ent-web-layout/internal/tools/clog"
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
