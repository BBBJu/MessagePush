package utils

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node // 雪花算法节点

func InitSnowflake(nodeID int64) error {
	if node != nil {
		return nil
	}
	var err error
	node, err = snowflake.NewNode(nodeID)
	return err
}

func GenerateSnowflakeID() string {
	return node.Generate().String()
}
