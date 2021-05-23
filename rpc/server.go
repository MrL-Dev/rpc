package rpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
)

// Server 结构体
type Server struct {
	network string
	address string
	conn    net.Conn                 // socket连接
	maps    map[string]reflect.Value // 函数字典
}

// NewServer 构造函数
func NewServer(network string, address string) *Server {
	return &Server{
		network: network,
		address: address,
		maps:    make(map[string]reflect.Value),
	}
}

// Register 注册函数
func (s *Server) Register(serviceImpl interface{}) {
	for i := 0; i < reflect.ValueOf(serviceImpl).NumMethod(); i++ {
		s.maps[reflect.TypeOf(serviceImpl).Method(i).Name] =
			reflect.ValueOf(serviceImpl).Method(i)
		log.Println("register " + reflect.TypeOf(serviceImpl).Method(i).Name)
	}
}

// Run 运行一个socket接收请求
func (s *Server) Run() {
	listen, err := net.Listen(s.network, s.address)
	if err != nil {
		panic(err)
	}
	log.Printf("listening at %s\n", s.address)
	for {
		s.conn, err = listen.Accept()
		if err != nil {
			continue
		}
		go s.handleConnect()
	}
}

func (s *Server) handleConnect() {
	for {
		transport := NewTransport(s.conn)
		req, err := transport.Receive()
		if err != nil {
			// 数据传输错误
			if err != io.EOF {
				log.Printf("read err: %v\n", err)
			}
			return
		}
		fn, exist := s.maps[req.Name]
		if !exist { // rpc接口不存在
			e := fmt.Sprintf("func %s does not exist", req.Name)
			log.Println(e)
			if err = transport.Send(&RpcData{Name: req.Name, Err: e}); err != nil {
				log.Printf("transport write err: %v\n", err)
			}
			continue
		}
		if len(req.Args)+1 != fn.Type().NumIn() {
			e := fmt.Sprintf("func %s does not exist", req.Name)
			log.Println(e)
			if err = transport.Send(&RpcData{Name: req.Name, Err: e}); err != nil {
				log.Printf("transport write err: %v\n", err)
			}
			continue
		}
		inArgs := make([]reflect.Value, 0, len(req.Args)+1)
		if !req.Deadline.IsZero() {
			ctx, cancel := context.WithDeadline(context.TODO(), req.Deadline)
			defer cancel()
			inArgs = append(inArgs, reflect.ValueOf(ctx))
		}
		for i := range req.Args {
			inArgs = append(inArgs, reflect.ValueOf(req.Args[i]))
		}
		out := fn.Call(inArgs)
		// 处理错误返回
		outArgs := make([]interface{}, len(out)-1)
		for i := 0; i < len(out)-1; i++ {
			outArgs[i] = out[i].Interface()
		}
		var e string
		if _, ok := out[len(out)-1].Interface().(error); !ok {
			e = ""
		} else {
			e = out[len(out)-1].Interface().(error).Error()
		}
		// 返回rsp
		err = transport.Send(&RpcData{Name: req.Name, Args: outArgs, Err: e})
		if err != nil {
			log.Printf("transport write err: %v\n", err)
		}
	}
}
