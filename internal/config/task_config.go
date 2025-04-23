package config

import (
    "errors"
)

// TaskConfig 适用于单个连接的配置
type TaskConfig struct {
    Name   string `toml:"name" json:"name" comment:"任务名称" `
    Sql    string `toml:"sql" json:"sql" comment:"SQL语句"`
    Status bool   `toml:"status" json:"status" comment:"是否启用"`
}

// TaskConfigMap 适用于多个连接的配置
type TaskConfigMap map[string]TaskConfig

var sqlConfig *TaskConfigMap

// GetSqlConfigAll 获取所有数据库配置
func GetSqlConfigAll() (*TaskConfigMap, error) {
    key := "task"

    if sqlConfig == nil {
        sqlConfig = &TaskConfigMap{}
    }

    // 如果配置不存在，则创建默认配置
    if !Exists(key) {
        err := SetTaskConfigAll(sqlConfig)
        if err != nil {
            return nil, err
        }
    }

    err := Unmarshal(key, sqlConfig)
    if err != nil {
        return nil, err
    }
    return sqlConfig, nil
}

// SetTaskConfigAll 设置所有数据库配置
func SetTaskConfigAll(conf *TaskConfigMap) error {
    key := "task"
    sqlConfig = conf
    return Marshal(key, conf)
}

// GetTaskConfig 获取指定数据库配置
func GetTaskConfig(name string) (*TaskConfig, error) {

    all, err := GetSqlConfigAll()
    if err != nil {
        return nil, err
    }

    allDb := *all

    // 判断all中是否存在
    if _, ok := allDb[name]; ok {
        db := allDb[name]
        return &db, nil
    }
    return nil, errors.New("sql config not found")
}

// SetTaskConfig 设置指定数据库配置
func SetTaskConfig(name string, conf *TaskConfig) error {

    all, err := GetSqlConfigAll()
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
