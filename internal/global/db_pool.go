package global

import (
    "app/internal/config"
    "app/internal/tools/safemap"
    "fmt"
    "golang.org/x/exp/slog"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "time"
)

type DbPoolType struct {
    DbName string
    Db     *gorm.DB
}

var DbPool struct {
    *safemap.Map[*DbPoolType]
}

// 使用自定义logger接管gorm日志
type dbLogWriter struct{}

func (w dbLogWriter) Printf(format string, args ...any) {
    slog.Warn(fmt.Sprintf(format, args...))
}

// ConnDb 开始链接数据库（初始化）
func ConnDb() error {

    // 连接商城数据库
    shopConf, err := config.GetDBConfig("shop")
    if err != nil {
        return err
    }

    gormConfig := &gorm.Config{
        Logger: logger.New(
            dbLogWriter{},
            logger.Config{
                SlowThreshold:             time.Duration(shopConf.SlowLogTime) * time.Millisecond, // 慢 SQL 阈值
                LogLevel:                  logger.Warn,                                            // 日志级别
                IgnoreRecordNotFoundError: true,                                                   // 忽略记录未找到错误
                Colorful:                  true,                                                   // 禁用彩色打印
            },
        ),
    }
    shopDb, err := gorm.Open(mysql.Open(shopConf.GenerateDsn()), gormConfig)
    if err != nil {
        return err
    }

    DbPool.Set("shop", &DbPoolType{
        DbName: "shop",
        Db:     shopDb,
    })

    // 连接erp数据库
    erpConf, err := config.GetDBConfig("erp")
    if err != nil {
        return err
    }

    erpConfig := &gorm.Config{
        Logger: logger.New(
            dbLogWriter{},
            logger.Config{
                SlowThreshold:             time.Duration(erpConf.SlowLogTime) * time.Millisecond, // 慢 SQL 阈值
                LogLevel:                  logger.Warn,                                           // 日志级别
                IgnoreRecordNotFoundError: true,                                                  // 忽略记录未找到错误
                Colorful:                  true,                                                  // 禁用彩色打印
            },
        ),
    }
    erpDb, err := gorm.Open(mysql.Open(erpConf.GenerateDsn()), erpConfig)
    if err != nil {
        return err
    }

    DbPool.Set("erp", &DbPoolType{
        DbName: "erp",
        Db:     erpDb,
    })

    return nil
}
