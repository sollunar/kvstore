package storage

import (
	"fmt"

	"github.com/tarantool/go-tarantool"
)

type TarantoolStorage struct {
	Conn *tarantool.Connection
}

func NewTarantoolStorage(host, port string) *TarantoolStorage {
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := tarantool.Connect(addr, tarantool.Opts{})
	if err != nil {
		panic("failed to connect to tarantool")
	}
	return &TarantoolStorage{Conn: conn}
}
