package service

import (
	"errors"
)

var (
	MyErrorSender = &ErrorSender{}
)

func InitSender() {
	InitLarkSender()
}

type MessageParams struct {
	ReceiveId string
	// 消息类型，支持：text、image等等
	MsgType string
	//针对飞书，传的是序列化后的json字符串，例如"{\"text\":\"不会哈气学哈气，机密咋摆你咋摆233\"}"
	Content string
}
type Sender interface {
	SendMessage(sp MessageParams) error
}

func GetSender(channel int) Sender {
	switch channel {
	case 1:
		//TODO: 没有复用连接， 之后改成连接池复用
		return MyLarkSender
	case 7788:
		return MyErrorSender
	}
	return nil
}

type ErrorSender struct {
}

func NewErrorSender() *ErrorSender {
	return &ErrorSender{}
}
func (e *ErrorSender) SendMessage(sp MessageParams) error {
	return errors.New("test retry")
}
