package database

import (
	"fmt"
	"github.com/owenzhou/ginrbac/support/facades"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func newDB() *gorm.DB {
	if facades.Config == nil{
		return nil
	}

	m := facades.Config.Mysql
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Dbname + "?charset=" + m.Charset + "&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	var loggerer logger.Interface
	//调试模式则打开sql输出
	if facades.Config.Debug {
		loggerer = logger.Default.LogMode(logger.Info)
	} else {
		loggerer = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //禁用表名复数
		},
		Logger: loggerer,
	})
	if err != nil {
		fmt.Println("Database error: can not connect to database.")
		return nil
	}
	return db
}
