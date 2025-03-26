package database

import (
	"fmt"
	"messagePush/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() error {
	// 构建 DSN 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MyConfig.MySQL.Username,
		config.MyConfig.MySQL.Pwd,
		config.MyConfig.MySQL.Host,
		config.MyConfig.MySQL.Port,
		config.MyConfig.MySQL.Database,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	fmt.Println("MySQL连接成功")
	return nil
}
