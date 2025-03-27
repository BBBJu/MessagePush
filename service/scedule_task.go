package service

import (
	"context"
	"fmt"
	"messagePush/database"
	"messagePush/models"
	"messagePush/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type CreateSceduleMessageParams struct {
	subject          string
	to               string
	channel          int
	sourceID         string
	templateId       int64
	templateData     string
	processTimeStamp int64 //处理时间最小单位为秒
}

func CreateSceduleMessage(params CreateSceduleMessageParams) {
	if params.processTimeStamp < time.Now().UnixNano()/1e9 {
		// 处理时间早于当前时间，直接发送
		fmt.Println("处理时间错误")
	}
	//TODO: 事务
	message := models.Message{
		SourceID:     params.sourceID,
		TemplateID:   params.templateId,
		TemplateData: params.templateData,
	}
	message.MsgID = utils.GenerateSnowflakeID()
	err := message.CreateMessage()
	if err != nil {
		fmt.Println(err.Error())
	}
	scheduleMessageQueue := models.ScheduleMessageQueue{
		MsgID:            message.MsgID,
		To:               params.to,
		Subject:          params.subject,
		Channel:          params.channel,
		ProcessTimeStamp: params.processTimeStamp,
		Status:           0,
	}
	//TODO: 保证数据库，Redis一致性
	err = scheduleMessageQueue.CreateScheduleMessageQueue()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = database.RedisClient.ZAdd(context.Background(), utils.Redis_Scedule_Task, redis.Z{
		Score:  float64(params.processTimeStamp),
		Member: message.MsgID,
	}).Err()
	if err != nil {
		fmt.Println(err.Error())
	}
}

var scanScript = redis.NewScript(`
-- KEYS[1]: 有序集合的key
-- ARGV[1]: 当前时间戳
local elements = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1])
if #elements > 0 then
    redis.call('ZREM', KEYS[1], unpack(elements))
end
return elements
`)

func StartRedisScanner() {
	if err := scanScript.Load(context.Background(), database.RedisClient).Err(); err != nil {
		panic(fmt.Sprintf("Failed to load Redis script: %v", err))
	}
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ProcessDelayedTasks()
			}
		}
	}()
}

func ProcessDelayedTasks() {
	ctx := context.Background()
	now := time.Now().Unix()

	keys := []string{utils.Redis_Scedule_Task}
	values := []interface{}{now}
	msgIDs := []string{}
	if result, err := scanScript.Run(ctx, database.RedisClient, keys, values...).Result(); err == nil {
		if elements, ok := result.([]interface{}); ok {
			for _, elem := range elements {
				if msgId, ok := elem.(string); ok {
					msgIDs = append(msgIDs, msgId)
				}
			}
		}
	} else {
		fmt.Println(err.Error())
	}
	processScheduledTasks(msgIDs)
}

func processScheduledTasks(messageIds []string) {
	if len(messageIds) == 0 {
		return
	}
	scheduleMessages, err := models.GetScheduleMessageQueuesByIds(messageIds)
	if err != nil {
		fmt.Println("获取任务详情失败:", err)
		return
	}
	messageQueue := make([]models.MessageQueue, len(scheduleMessages))
	for i, scheduleMessage := range scheduleMessages {
		messageQueue[i] = models.MessageQueue{
			MsgID:      scheduleMessage.MsgID,
			To:         scheduleMessage.To,
			Subject:    scheduleMessage.Subject,
			Channel:    scheduleMessage.Channel,
			Status:     MessageStatusCreated,
			Priority:   VIPPriority,
			OrderBy:    time.Now().Unix() - VIPPriority,
			RetryCount: 0,
		}
		scheduleMessages[i].Status = 1
	}
	err = models.BatchCreateMessageQueue(messageQueue)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = models.BatchUpdateScheduleMessageQueueStatus(scheduleMessages)
	if err != nil {
		fmt.Println(err.Error())
	}
}
