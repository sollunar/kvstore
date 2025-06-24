package kvstore

import (
	"fmt"

	"github.com/sollunar/kvstore-api/pkg/storage"
	"github.com/tarantool/go-tarantool/v2"
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
	data, err := s.Conn.Do(
		tarantool.NewSelectRequest("kvstore").Key([]interface{}{key}),
	).Get()

	if err != nil || len(data) == 0 {
		return "", ErrKeyNotFound
	}

	return data[0].([]interface{})[1].(string), nil
}

func (s *KVStore) Set(key, value string) error {
	_, err := s.Conn.Do(
		tarantool.NewUpsertRequest("kvstore").Tuple([]interface{}{key, value}).Operations(tarantool.NewOperations().Assign(1, value)),
	).Get()

	fmt.Println(err.Error())

	return err
}

func (s *KVStore) Delete(key string) error {

	data, err := s.Conn.Do(
		tarantool.NewDeleteRequest("kvstore").Key([]interface{}{key}),
	).Get()

	if err != nil || len(data) == 0 {
		return ErrKeyNotFound
	}

	return nil
}
