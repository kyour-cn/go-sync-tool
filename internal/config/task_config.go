package config

import (
    "context"
    "errors"
    "github.com/go-gourd/gourd/event"
)

// TaskConfig 适用于单个连接的配置
type TaskConfig struct {
    Name         string `toml:"name" json:"name" comment:"任务名称" `
    Sql          string `toml:"sql" json:"sql" comment:"SQL语句"`
    IntervalTime int    `toml:"interval_time" json:"interval_time" comment:"间隔时间-秒"`
    Status       bool   `toml:"status" json:"status" comment:"是否启用"`
}

// TaskConfigMap 适用于多个连接的配置
type TaskConfigMap map[string]TaskConfig

var taskConfig *TaskConfigMap

// GetTaskConfigAll 获取所有数据库配置
func GetTaskConfigAll() (*TaskConfigMap, error) {
    key := "task"

    if taskConfig == nil {
        taskConfig = &TaskConfigMap{}
    }

    // 如果配置不存在，则创建默认配置
    if !Exists(key) {
        err := SetTaskConfigAll(taskConfig)
        if err != nil {
            return nil, err
        }
    }

    err := Unmarshal(key, taskConfig)
    if err != nil {
        return nil, err
    }
    return taskConfig, nil
}

// SetTaskConfigAll 设置所有数据库配置
func SetTaskConfigAll(conf *TaskConfigMap) error {
    key := "task"
    taskConfig = conf

    defer event.Trigger("task.config", context.Background())

    return Marshal(key, conf)
}

// GetTaskConfig 获取指定数据库配置
func GetTaskConfig(name string) (*TaskConfig, error) {

    all, err := GetTaskConfigAll()
    if err != nil {
        return nil, err
    }

    allDb := *all

    // 判断all中是否存在
    if _, ok := allDb[name]; ok {
        db := allDb[name]
        return &db, nil
    }
    return nil, errors.New("task config not found")
}

// SetTaskConfig 设置指定数据库配置
func SetTaskConfig(name string, conf *TaskConfig) error {

    all, err := GetTaskConfigAll()
    if err != nil {
        return err
    }

    allDb := *all

    // 判断all中是否存在
    if _, ok := allDb[name]; ok {
        allDb[name] = *conf
    } else {
        allDb[name] = *conf
    }

    return SetTaskConfigAll(&allDb)
}
