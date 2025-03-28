package main

import (
	"fmt"
	"log"
	"messagePush/config"
	"messagePush/database"
	"messagePush/models"
	"messagePush/service"
	"messagePush/utils"
	"os"
	"time"

	"github.com/withlin/canal-go/client"
	pbe "github.com/withlin/canal-go/protocol/entry"
	"google.golang.org/protobuf/proto"
)

func main() {
	//mySQL连接
	config.InitConfig()
	database.InitMySQL()
	database.InitRedis()
	service.InitSender()
	models.Migrate()
	utils.InitSnowflake(0)

	// 启动数据库轮询
	go PollDatabase()
	// 新增Redis扫描器
	go service.StartRedisScanner()
	//压力测试
	go service.StartStressTest()

	// 保持主程序运行
	select {}
}

func PollDatabase() {
	ticker := time.NewTicker(1 * time.Second) // 每1秒轮询一次
	defer ticker.Stop()
	for {
		fmt.Println("开始轮询数据库")
		select {
		case <-ticker.C:
			messageQueues, err := models.GetPendingMessages(1000) // 默认一次获取1000条消息
			if err != nil {
				log.Println("Error fetching messages:", err)
				continue
			}
			if len(messageQueues) > 0 {
				service.HandleMessage(messageQueues) // 处理消息
			}
		}
	}
}

func Canal() {
	// 192.168.199.17 替换成你的canal server的地址
	// example 替换成-e canal.destinations=example 你自己定义的名字
	connector := client.NewSimpleCanalConnector("127.0.0.1", 11111, "", "", "example", 60000, 60*60*1000)
	err := connector.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//常见例子：
	//
	//  1.  所有表：.*   or  .*\\..*
	//	2.  canal schema下所有表： canal\\..*
	//	3.  canal下的以canal打头的表：canal\\.canal.*
	//	4.  canal schema下的一张表：canal\\.test1
	//  5.  多个规则组合使用：canal\\..*,mysql.test1,mysql.test2 (逗号分隔)

	err = connector.Subscribe("MessagePush\\.message")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		message, err := connector.Get(100, nil, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(1000 * time.Millisecond)
			fmt.Println("===没有数据了===")
			continue
		}

		//printEntry(message.Entries)
		for _, entry := range message.Entries {
			if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
				continue
			}
			rowChange := new(pbe.RowChange)

			err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
			checkError(err)
			if rowChange != nil {
				eventType := rowChange.GetEventType()
				header := entry.GetHeader()
				fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))
				msgIds := []string{}
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						msgId := getColumnValue(rowData.GetAfterColumns(), "msg_id")
						msgIds = append(msgIds, msgId)
						printColumn(rowData.GetAfterColumns())
					}
				}
				//: 修改参数
				service.HandleMessage(nil)
			}

		}
	}
}

func getColumnValue(columns []*pbe.Column, name string) string {
	for _, col := range columns {
		if col.GetName() == name {
			return col.GetValue()
		}
	}
	return ""
}

func printEntry(entrys []pbe.Entry) {

	for _, entry := range entrys {
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(pbe.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)
		if rowChange != nil {
			eventType := rowChange.GetEventType()
			header := entry.GetHeader()
			fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))

			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_DELETE {
					printColumn(rowData.GetBeforeColumns())
				} else if eventType == pbe.EventType_INSERT {
					printColumn(rowData.GetAfterColumns())
				} else {
					fmt.Println("-------> before")
					printColumn(rowData.GetBeforeColumns())
					fmt.Println("-------> after")
					printColumn(rowData.GetAfterColumns())
				}
			}
		}
	}
}

func printColumn(columns []*pbe.Column) {
	for _, col := range columns {
		fmt.Println(fmt.Sprintf("%s : %s  update= %t", col.GetName(), col.GetValue(), col.GetUpdated()))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
