package helloworld

import (
	"context"
	"encoding/gob"
	"errors"
	"log"
)

func init() {
	gob.Register(GreetReq{})
	gob.Register(GreetRsp{})
}

type GreetReq struct {
	Name string `json:"name"`
	Echo string `json:"echo"`
}
type GreetRsp struct {
	Msg  string `json:"msg"`
	Echo string `json:"echo"`
}

func (s *HelloWorldServer) Greet(ctx context.Context, req GreetReq) (GreetRsp, error) {
	if req.Name == "" {
		return GreetRsp{}, errors.New("invalid param")
	}
	log.Printf("req=%+v,ctx=%+v\n", req, ctx)
	return GreetRsp{Msg: "Hi," + req.Name, Echo: req.Echo}, nil
}

type HelloWorldServer struct {
}
