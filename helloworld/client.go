package helloworld

import (
	"context"
	"myrpc/rpc"
)

type HelloWorldClient struct {
	cli *rpc.Client
}

func (c *HelloWorldClient) Greet(ctx context.Context, req *GreetReq) (*GreetRsp, error) {
	var r = rpc.RpcData{
		Name: "Greet",
		Args: []interface{}{*req},
	}
	d, ok := ctx.Deadline()
	if ok {
		r.Deadline = d
	}
	rsp, err := c.cli.Call(&r)
	if err != nil {
		return nil, err
	}
	p := &GreetRsp{}
	if len(rsp.Args) != 0 {
		*p = rsp.Args[0].(GreetRsp)
	}
	return p, nil
}

func NewHelloWorldClient(network string, address string) (*HelloWorldClient, error) {
	var err error
	client := &HelloWorldClient{}
	client.cli, err = rpc.NewClient(network, address)
	if err != nil {
		return nil, err
	}
	return client, nil
}
