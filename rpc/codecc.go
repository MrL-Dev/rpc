package rpc

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type ICodec interface {
	Decode(data []byte) (*RpcData, error)
	Encode(data *RpcData) ([]byte, error)
}

type JsonCodec struct {
}

func (c *JsonCodec) Decode(data []byte) (*RpcData, error) {
	dst := &RpcData{}
	if err := json.Unmarshal(data, dst); err != nil {
		return nil, err
	}
	return dst, nil
}

func (c *JsonCodec) Encode(data *RpcData) ([]byte, error) {
	dst, err := json.Marshal(data)
	if err != nil {
		return []byte{}, err
	}
	return dst, nil
}

type GobCodec struct {
}

func (c *GobCodec) Decode(data []byte) (*RpcData, error) {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	dst := &RpcData{}
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return dst, nil
}

func (c *GobCodec) Encode(data *RpcData) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
