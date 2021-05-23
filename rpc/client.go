package rpc

import (
	"errors"
	"log"
	"net"
)

// Client 结构体
type Client struct {
	conn net.Conn // socket连接
}

// NewClient 构造函数
func NewClient(network string, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Call(data *RpcData) (*RpcData, error) {
	transport := NewTransport(c.conn)
	if err := transport.Send(data); err != nil {
		// 数据传输错误
		log.Printf("read err: %v\n", err)
		return nil, err
	}
	// 返回rsp
	rsp, err := transport.Receive()
	if err != nil {
		return nil, err
	}
	if rsp.Err != "" {
		return nil, errors.New(rsp.Err)
	}
	return rsp, nil
}
