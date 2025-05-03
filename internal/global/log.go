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

    //判断日志仅保留最后100条
    if len(Logs) > 100 {
        Logs = Logs[len(Logs)-100:]
    }
    //Logs = append(Logs, Log{
    //    Time:    time.Now(),
    //    Level:   "info",
    //    Message: string(p),
    //})

    // 往头部插入数据
    Logs = append([]Log{{
        Time:    time.Now(),
        Level:   l,
        Message: msg,
    }}, Logs...)

    // 刷新日志UI
    event.Trigger("window.invalidate", context.Background())

}
