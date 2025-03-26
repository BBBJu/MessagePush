package service

import (
	"fmt"
	"math/rand"
	"messagePush/models"
	"messagePush/utils"
	"time"
)

const (
	MessageStatusCreated = 1
	MessageStatusSending = 2
	MessageStatusSuccess = 3
	MessageStatusFail    = 4
)

type CreateMessageParams struct {
	subject  string
	to       string
	channel  int
	sourceID string
}

func CreateMessage(params CreateMessageParams) {
	message := models.Message{
		Subject:  params.subject,
		To:       params.to,
		Channel:  params.channel,
		SourceID: params.sourceID,
	}
	message.MsgID = utils.GenerateSnowflakeID()
	err := message.CreateMessage()
	if err != nil {
		fmt.Println(err.Error())
	}

	messageQueue := models.MessageQueue{
		Subject:  params.subject,
		To:       params.to,
		Channel:  params.channel,
		MsgID:    message.MsgID,
		Status:   MessageStatusCreated,
		Priority: 1,
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
func HandleMessage(msgIds []string) {
	messageQueues, err := models.GetMessageQueueByMsgIDs(msgIds)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//TODO: 先写成循环更新， 后面可以考虑批量更新以及多携程
	for _, messageQueue := range messageQueues {
		if messageQueue.Status == MessageStatusCreated {
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
			//TODO: 修改成template的形式
			println(message.TemplateData)
			content := fmt.Sprintf("{\"text\":\"不会哈气学哈气，机密咋摆你咋摆%d\"}", RandInt(0, 100))
			messageParams := MessageParams{
				ReceiveId: messageQueue.To,
				Content:   content,
				//TODO 暂时只支持text类型
				MsgType: "text",
			}
			err = SendMessage(GetSender(messageQueue.Channel), messageParams)
			if err != nil {
				fmt.Println(err.Error())
				messageQueue.Status = MessageStatusFail
				//TODO: 失败重试
				err = messageQueue.UpdateMessageQueue()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			} else {
				messageQueue.Status = MessageStatusSuccess
				err = messageQueue.UpdateMessageQueue()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		} else {
			fmt.Println("消息状态不是created， 不处理")
		}
	}
}

func GetMessageFromCanal() []string {
	return nil
}

// 在包初始化时设置随机种子
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 新增随机数生成函数
func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
