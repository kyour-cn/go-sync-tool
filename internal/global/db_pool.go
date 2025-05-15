package global

import (
	"app/internal/config"
	"app/internal/orm/shop_query"
	"app/internal/tools/safemap"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"net/url"
	"time"
)

var DbPool = safemap.New[*gorm.DB]()

// 使用自定义logger接管gorm日志
type dbLogWriter struct{}

func (w dbLogWriter) Printf(format string, args ...any) {
	slog.Debug(fmt.Sprintf(format, args...))
}

// ConnDb 开始链接数据库（初始化）
func ConnDb() error {

	if DbPool.Len() > 0 {
		// 判断数据库连接是否正常
		for _, v := range *DbPool.GetMap() {
			db, err := v.DB()
			if err != nil {
				return err
			}
			err = db.Ping()
			if err != nil {
				return err
			}
		}
		return nil
	}

	// 连接商城数据库
	shopConf, err := config.GetDBConfig("shop")
	if err != nil {
		return err
	}

	shopConf.Param = "parseTime=True&timeout=30s&readTimeout=60s&writeTimeout=60s"

	gormConfig := &gorm.Config{
		Logger: logger.New(
			dbLogWriter{},
			logger.Config{
				SlowThreshold:             time.Duration(shopConf.SlowLogTime) * time.Millisecond, // 慢 SQL 阈值
				LogLevel:                  logger.Warn,                                            // 日志级别
				IgnoreRecordNotFoundError: true,                                                   // 忽略记录未找到错误
				Colorful:                  false,                                                  // 彩色打印
			},
		),
		PrepareStmt: true, //缓存预编译语句
	}
	shopDb, err := gorm.Open(mysql.Open(shopConf.GenerateDsn()), gormConfig)
	if err != nil {
		return err
	}

	shopDbConn, _ := shopDb.DB()
	// 设置空闲连接池中的最大连接数。
	shopDbConn.SetMaxIdleConns(20)
	// 设置数据库的最大打开连接数。
	shopDbConn.SetMaxOpenConns(100)

	// 设置shop查询默认的数据库连接
	shop_query.SetDefault(shopDb)

	DbPool.Set("shop", shopDb)

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
		PrepareStmt: true, //缓存预编译语句
	}

	var erpDb *gorm.DB
	if erpConf.Type == "sqlserver" {
		// 首先建立ODBC连接
		var erpOdbc *sql.DB
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
		// 然后将连接传递给GORM
		erpDb, err = gorm.Open(sqlserver.New(sqlserver.Config{
			Conn: erpOdbc,
		}), erpConfig)
		if err != nil {
			return err
		}
	} else if erpConf.Type == "mysql" {
		erpDb, err = gorm.Open(mysql.Open(erpConf.GenerateDsn()), erpConfig)
	} else {
		// TODO: 其它数据库类型待实现
		return errors.New("ODBC暂未支持的数据库类型：" + erpConf.Type)
	}

	DbPool.Set("erp", erpDb)

	return nil
}

func CloseDb() error {
	// 遍历数据库连接池，关闭每个连接
	dbMap := DbPool.GetMap()
	if dbMap == nil {
		return nil
	}
	for _, v := range *dbMap {
		db, err := v.DB()
		if err != nil {
			return err
		}
		err = db.Close()
		if err != nil {
			return err
		}
	}

	DbPool.Clear()

	return nil
}
