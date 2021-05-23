package rpc

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(RpcData{})
}

// RpcData ...
type RpcData struct {
	Name     string        // 函数名
	Args     []interface{} // 参数
	Err      string        // 错误
	Deadline time.Time
}
