package kvstore

import (
	"github.com/sollunar/kvstore-api/pkg/storage"
	"github.com/tarantool/go-tarantool"
)

type KVStore struct {
	*storage.TarantoolStorage
}

func NewKVRepository(storage *storage.TarantoolStorage) *KVStore {
	return &KVStore{
		TarantoolStorage: storage,
	}
}

func (s *KVStore) Get(key string) (string, error) {
	resp, err := s.Conn.Select("kv", "primary", 0, 1, tarantool.IterEq, []any{key})
	if err != nil || len(resp.Data) == 0 {
		return "", ErrKeyNotFound
	}
	return resp.Data[0].([]any)[1].(string), nil
}

func (s *KVStore) Set(key, value string) error {
	_, err := s.Conn.Replace("kv", []any{key, value})
	return err
}

func (s *KVStore) Delete(key string) error {
	resp, err := s.Conn.Delete("kv", "primary", []any{key})
	if err != nil || len(resp.Data) == 0 {
		return ErrKeyNotFound
	}
	return nil
}
