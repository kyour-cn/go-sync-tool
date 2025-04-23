package domain

import (
    "time"
)

type Log struct {
    Time    time.Time
    Level   string
    Message string
}

var Logs []Log

// WriteConsoleLog 写入日志到控制台
func WriteConsoleLog(p []byte) (n int, err error) {

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
        Level:   "info",
        Message: string(p),
    }}, Logs...)

    return len(p), nil
}
