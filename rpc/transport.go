package rpc

import (
	"encoding/binary"
	"io"
	"net"
)

// Transport struct
type Transport struct {
	conn  net.Conn
	codec ICodec
}

// NewTransport creates a transport
func NewTransport(conn net.Conn) *Transport {
	return &Transport{
		conn: conn,
		// codec: &JsonCodec{},
		codec: &GobCodec{},
	}
}

// Send data
func (t *Transport) Send(req *RpcData) error {
	b, err := t.codec.Encode(req)
	if err != nil {
		return err
	}
	buf := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(b)))
	copy(buf[4:], b)
	_, err = t.conn.Write(buf)
	return err
}

// Receive data
func (t *Transport) Receive() (*RpcData, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(t.conn, header)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(t.conn, data)
	if err != nil {
		return nil, err
	}
	rsp, err := t.codec.Decode(data)
	return rsp, err
}
