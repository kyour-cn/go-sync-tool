package global

import (
	"context"
	"github.com/go-gourd/gourd/event"
	"time"
)

type Log struct {
	Time    time.Time
	Level   string
	Message string
}

var Logs []Log

// WriteConsoleLog 写入日志到控制台
func WriteConsoleLog(l string, msg string) {

	// 往头部插入新日志
	Logs = append([]Log{{
		Time:    time.Now(),
		Level:   l,
		Message: msg,
	}}, Logs...)

	// 截断保留前100条（新插入的日志在头部）
	if len(Logs) > 100 {
		Logs = Logs[:100]
	}

	// 刷新日志UI
	event.Trigger("window.invalidate", context.Background())
}
