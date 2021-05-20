package rpc

import (
	"net"
	"reflect"
)

// server 结构体
type server struct {
	conn net.Conn                 //socket连接
	maps map[string]reflect.Value //函数字典
}

// NewServer 构造函数
func NewServer() *server {
	return &server{
		maps: make(map[string]reflect.Value),
	}
}

// Register 注册函数
func (s *server) Register(fname string, fun interface{}) {
	if _, ok := s.maps[fname]; !ok {
		s.maps[fname] = reflect.ValueOf(fun)
	}
}

// Run 运行一个socket接收请求
func (s *server) Run() {
	listen, err := net.Listen("tcp4", ":3001")
	if err != nil {
		panic(err)
	}
	for {
		s.conn, err = listen.Accept()
		if err != nil {
			continue
		}
		go s.handleConnect()
	}
}

func (s *server) handleConnect() {
	for {
		// header := make([]byte, 4)
		// if _, err := s.conn.Read(header); err != nil {
		// 	continue
		// }
		// bodyLen := binary.BigEndian.Uint32(header)
		// body := make([]byte, int(bodyLen))
		// if _, err := s.conn.Read(body); err != nil {
		// 	continue
		// }

		transport := NewTransport(s.conn)
		req, err := transport.Receive()
		if err != nil {
			continue
		}
		inArgs := make([]reflect.Value, len(req.Args))
		for i := range req.Args {
			inArgs[i] = reflect.ValueOf(req.Args[i])
		}
		fn, exist := s.maps[req.Name]
		if !exist {
			continue
		}
		fn.Call(inArgs)
		// package error argument
		// var e string
		// if _, ok := out[len(out)-1].Interface().(error); !ok {
		// 	e = ""
		// } else {
		// 	e = out[len(out)-1].Interface().(error).Error()
		// }
	}
}
