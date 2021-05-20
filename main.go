package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"myrpc/rpc"
	"net"
)

func main() {
	fmt.Println("my rpc demo")
	go runServer()
	go runClient()
	select {}
}

func runClient() {
	var req = rpc.RpcData{
		Name: "fn",
		Args: []interface{}{1, "aaa"},
	}
	rpcCall(req)
}

func rpcCall(data rpc.RpcData) {
	conn, err := net.Dial("tcp4", "127.0.0.1:3001")
	if err != nil {
		panic(err)
	}
	req, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 4+len(req))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(req)))
	copy(buf[4:], req)
	_, err = conn.Write(buf)
	if err != nil {
		panic(err)
	}
}

func runServer() {
	srv := rpc.NewServer()
	srv.Register("fn", fn)
	srv.Run()
}

func fn(args ...interface{}) {
	fmt.Println(args...)
}
