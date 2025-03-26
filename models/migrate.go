package models

import (
	"messagePush/database"
)

// 迁移
func Migrate() {
	database.DB.AutoMigrate(&Message{})
	database.DB.AutoMigrate(&MessageQueue{})
}
