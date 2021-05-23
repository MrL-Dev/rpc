package main

import (
	"context"
	"fmt"
	"log"
	"myrpc/helloworld"
	"myrpc/rpc"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("my rpc demo test...")
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
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*1)
	defer cancel()
	rsp, err := cli.Greet(ctx, &helloworld.GreetReq{Name: "nobody", Echo: "lalala"})
	log.Printf("err=%+v,rsp=%+v\n", err, rsp)

	rsp, err = cli.Greet(ctx, &helloworld.GreetReq{Name: ""})
	log.Printf("err=%+v,rsp=%+v\n", err, rsp)
}
