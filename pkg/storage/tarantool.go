package storage

import (
	"context"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

type TarantoolStorage struct {
	Conn *tarantool.Connection
}

func NewTarantoolStorage(host, port string) *TarantoolStorage {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address: host + ":" + port,
	}

	opts := tarantool.Opts{
		Timeout: time.Second,

	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		panic("failed to connect to tarantool")
	}
	return &TarantoolStorage{Conn: conn}
}
