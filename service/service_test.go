package service

import (
	"fmt"
	"messagePush/config"
	"messagePush/database"
	"messagePush/models"
	"messagePush/utils"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	Init()
	params := CreateMessageParams{
		subject:  "Test Subject",
		to:       utils.ReceiveId,
		channel:  1,
		sourceID: "Test Source",
	}
	for i := 0; i < 50; i++ {
		CreateMessage(params)
	}
}

func TestGetMessage(t *testing.T) {
	Init()
	msgID := []string{"1904538274127413248"}
	HandleMessage(msgID)
}

func TestP(t *testing.T) {
	fmt.Sprintf("{\"text\":\"不会哈气学哈气，机密咋摆你咋摆233\"}")
}
func Init() {
	config.InitConfig()
	database.InitMySQL()
	models.Migrate()
	utils.InitSnowflake(0)
}
