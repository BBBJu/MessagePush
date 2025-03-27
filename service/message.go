package service

import (
	"encoding/json"
	"fmt"
	"messagePush/models"
	"messagePush/utils"
	"time"
)

const (
	MessageStatusCreated = 1
	MessageStatusSending = 2
	MessageStatusSuccess = 3
	MessageStatusFail    = 4
	MessageStatusDead    = 5 //超过重试次数，代表死亡
)

const (
	NormalPriority = 1
	VIPPriority    = 100
)

type CreateMessageParams struct {
	subject      string
	to           string
	channel      int
	sourceID     string
	templateId   int64
	templateData string
	priority     int64
}

func CreateMessage(params CreateMessageParams) {
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
	messageQueue := models.MessageQueue{
		MsgID:      message.MsgID,
		To:         params.to,
		Subject:    params.subject,
		Channel:    params.channel,
		Status:     MessageStatusCreated,
		Priority:   params.priority,
		OrderBy:    time.Now().Unix() - params.priority,
		RetryCount: 0,
	}
	err = messageQueue.CreateMessageQueue()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SendMessage(s Sender, messageParams MessageParams) error {
	if s == nil {
		fmt.Println("sender is nil, 表示输入的channel有误")
		return fmt.Errorf("sender is nil")
	}
	err := s.SendMessage(messageParams)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// 从数据库MessageQueue货期消息，并发送
func HandleMessage(messageQueues []models.MessageQueue) {
	//TODO: 先写成循环更新， 后面可以考虑批量更新以及多携程
	for _, messageQueue := range messageQueues {
		if messageQueue.Status == MessageStatusCreated || messageQueue.Status == MessageStatusFail {
			messageQueue.Status = MessageStatusSending
			err := messageQueue.UpdateMessageQueue()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			message, err := models.GetMessageByMsgId(messageQueue.MsgID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			content := DoTemplate(&message)
			messageParams := MessageParams{
				ReceiveId: messageQueue.To,
				Content:   content,
				//TODO 暂时只支持text类型
				MsgType: "text",
			}
			err = SendMessage(GetSender(messageQueue.Channel), messageParams)
			if err != nil {
				//最多重试五次
				if messageQueue.RetryCount < 5 {
					messageQueue.RetryCount += 1
					messageQueue.Status = MessageStatusFail
					fmt.Println("重试", messageQueue.MsgID, "次数", messageQueue.RetryCount)
					//重试不修改order_by优先级，追求时效性
					messageQueue.UpdateMessageQueue()
				} else {
					messageQueue.Status = MessageStatusDead
					messageQueue.UpdateMessageQueue()
				}
			} else {
				messageQueue.Status = MessageStatusSuccess
				messageQueue.UpdateMessageQueue()
			}
		} else {
			fmt.Println("消息状态不是created， 不处理")
		}
	}
}

func MessageDaemon() {
	ticker := time.NewTicker(10 * time.Second) // 创建10秒间隔的定时器
	defer ticker.Stop()                        // 函数退出时停止ticker
	for {
		select {
		case <-ticker.C: // 每10秒触发一次

			messageQueues, err := models.GetFailedMessageQueue()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			msgIds := make([]string, len(messageQueues))
			for i, messageQueue := range messageQueues {
				msgIds[i] = messageQueue.MsgID
			}
			//TODO: 修改参数
			HandleMessage(nil)
		}
	}
}

func DoTemplate(message *models.Message) string {
	//默认使用第一个模板
	if message.TemplateID == 0 {
		message.TemplateID = 1
	}
	var data map[string]interface{}
	err := json.Unmarshal([]byte(message.TemplateData), &data)
	if err != nil {
		fmt.Println(err.Error())
	}
	myTemlate := models.Template{
		ID: message.TemplateID,
	}
	myTemlate.GetTemplateById()

	content := utils.GetContentAfterTemplate(data, myTemlate)
	content = fmt.Sprintf("{\"text\":\"%s\"}", content)
	return content
}
