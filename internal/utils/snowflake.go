package utils

import (
	"github.com/bwmarrin/snowflake"
)

var sf *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	sf = node
}

func GenerateID() string {
	return sf.Generate().String()
}
