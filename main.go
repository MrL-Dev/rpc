package main

import (
	"context"
	"fmt"
	"myrpc/helloworld"
	xlog "myrpc/log"
	"myrpc/rpc"
	"time"
)

func main() {
	fmt.Println("my rpc demo test...")
	xlog.Info("aaa")
	go runHelloServer()
	time.Sleep(time.Second)
	go runHelloClient()

	select {}

}

func runHelloServer() {
	s := rpc.NewServer("tcp4", "127.0.0.1:3002")
	s.Register(&helloworld.HelloWorldServer{})
	s.Run()
}

func runHelloClient() {
	cli, err := helloworld.NewHelloWorldClient("tcp4", "127.0.0.1:3002")
	if err != nil {
		xlog.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*1)
	defer cancel()
	rsp, err := cli.Greet(ctx, &helloworld.GreetReq{Name: "nobody", Echo: "lalala"})
	xlog.Infof("err=%+v,rsp=%+v\n", err, rsp)

	rsp, err = cli.Greet(ctx, &helloworld.GreetReq{Name: ""})
	xlog.Infof("err=%+v,rsp=%+v\n", err, rsp)
}
