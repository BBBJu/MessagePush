package service

import (
	"fmt"
	"messagePush/config"
	"messagePush/database"
	"messagePush/models"
	"messagePush/utils"
	"testing"
	"time"
)

func TestCreateMessage(t *testing.T) {
	Init()
	params1 := CreateMessageParams{
		subject:      "Test Subject",
		to:           utils.ReceiveIdYang,
		channel:      7788,
		sourceID:     "Test Source",
		templateId:   1,
		templateData: "{\"username\": \"YANG\"}",
		priority:     10,
	}
	for i := 0; i < 100; i++ {
		CreateMessage(params1)
	}
	params2 := CreateMessageParams{
		subject:      "Test Subject",
		to:           utils.ReceiveIdYang,
		channel:      1,
		sourceID:     "Test Source",
		templateId:   4,
		templateData: "{\"username\": \"YANG\"}",
		priority:     1000,
	}
	for i := 0; i < 100; i++ {
		CreateMessage(params2)
	}
}
func TestGetTemplate(t *testing.T) {
	Init()
	myTemlate := models.Template{
		Name: "lark_hajimi",
	}
	err := myTemlate.GetTemplateByName()
	if err != nil {
		t.Fatalf("获取模板失败: %v", err)
	}

	// 准备模板数据
	data := map[string]interface{}{
		"username": "测试用户",
	}
	utils.GetContentAfterTemplate(data, myTemlate)
}

func TestCreateScheduleMessage(t *testing.T) {
	Init()
	params1 := CreateSceduleMessageParams{
		subject:          "Test Subject",
		to:               utils.ReceiveIdYang,
		channel:          1,
		sourceID:         "Test Source",
		templateId:       4,
		templateData:     "{\"username\": \"YANG\"}",
		processTimeStamp: time.Now().UnixNano()/1e9 + 10,
	}
	CreateSceduleMessage(params1)
}

func TestP(t *testing.T) {
	Init()
	now := time.Now().Unix()
	fmt.Println(now)
}
func Init() {
	config.InitConfig()
	database.InitMySQL()
	database.InitRedis()
	models.Migrate()
	utils.InitSnowflake(0)
	InitSender()
}
