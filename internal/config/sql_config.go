package config

import (
    "errors"
)

// SqlConfig 适用于单个连接的配置
type SqlConfig struct {
    Name string `toml:"name" json:"name" comment:"任务名称" `
    Sql  string `toml:"sql" json:"sql" comment:"SQL语句"`
}

// SqlConfigMap 适用于多个连接的配置
type SqlConfigMap map[string]SqlConfig

var sqlConfig *SqlConfigMap

// GetSqlConfigAll 获取所有数据库配置
func GetSqlConfigAll() (*SqlConfigMap, error) {
    key := "sql"

    if sqlConfig == nil {
        sqlConfig = &SqlConfigMap{}
    }

    // 如果配置不存在，则创建默认配置
    if !Exists(key) {
        err := SetSqlConfigAll(sqlConfig)
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

// SetSqlConfigAll 设置所有数据库配置
func SetSqlConfigAll(conf *SqlConfigMap) error {
    key := "sql"
    sqlConfig = conf
    return Marshal(key, conf)
}

// GetSqlConfig 获取指定数据库配置
func GetSqlConfig(name string) (*SqlConfig, error) {

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

// SetSqlConfig 设置指定数据库配置
func SetSqlConfig(name string, conf *SqlConfig) error {

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

    return SetSqlConfigAll(&allDb)
}
