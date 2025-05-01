package global

import (
    "app/internal/config"
    "app/internal/tools/safemap"
    "database/sql"
    "errors"
    "fmt"
    _ "github.com/alexbrainman/odbc"
    "golang.org/x/exp/slog"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlserver"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "net/url"
    "time"
)

type DbPoolType struct {
    DbName string
    Db     *gorm.DB
}

var DbPool *safemap.Map[*DbPoolType]

// 使用自定义logger接管gorm日志
type dbLogWriter struct{}

func (w dbLogWriter) Printf(format string, args ...any) {
    slog.Warn(fmt.Sprintf(format, args...))
}

// ConnDb 开始链接数据库（初始化）
func ConnDb() error {

    if DbPool == nil {
        DbPool = safemap.New[*DbPoolType]()
    }

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

    // 首先建立ODBC连接
    var erpOdbc *sql.DB
    if erpConf.Type == "sqlserver" {
        dsn := fmt.Sprintf("DRIVER={SQL Server};SERVER=%s,%d;DATABASE=%s;UID=%s;PWD=%s",
            erpConf.Host,
            erpConf.Port,
            erpConf.Database,
            erpConf.User,
            url.QueryEscape(erpConf.Pass),
        )
        // 创建数据库连接
        erpOdbc, err = sql.Open("odbc", dsn)
        if err != nil {
            return err
        }
    } else {
        // TODO: 其它数据库类型待实现
        return errors.New("ODBC暂不支持的数据库类型：" + erpConf.Type)
    }

    // 然后将连接传递给GORM
    erpDb, err := gorm.Open(sqlserver.New(sqlserver.Config{
        Conn: erpOdbc,
    }), erpConfig)
    if err != nil {
        return err
    }

    DbPool.Set("erp", &DbPoolType{
        DbName: "erp",
        Db:     erpDb,
    })

    return nil
}

func CloseDb() error {

    // 遍历数据库连接池，关闭每个连接
    for _, v := range DbPool.Values() {
        db, err := v.Db.DB()
        if err != nil {
            return err
        }
        err = db.Close()
        if err != nil {
            return err
        }
    }

    return nil
}
