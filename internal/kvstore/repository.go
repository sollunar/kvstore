package kvstore

import (
	"github.com/sollunar/kvstore-api/pkg/storage"
	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"
)

type KVStore struct {
	*storage.TarantoolStorage
	log *zap.Logger
}

func NewKVRepository(storage *storage.TarantoolStorage, logger *zap.Logger) *KVStore {
	return &KVStore{
		TarantoolStorage: storage,
		log:              logger,
	}
}

func (s *KVStore) Get(key string) (string, error) {
	s.log.Info("GET request", zap.String("key", key))

	data, err := s.Conn.Do(
		tarantool.NewSelectRequest("kvstore").Key([]interface{}{key}),
	).Get()

	if err != nil {
		s.log.Error("GET failed", zap.String("key", key), zap.Error(err))
		return "", ErrInternal
	}

	if len(data) == 0 {
		s.log.Warn("GET key not found", zap.String("key", key))
		return "", ErrKeyNotFound
	}

	val, ok := data[0].([]interface{})[1].(string)
	if !ok {
		s.log.Error("GET value type mismatch", zap.String("key", key))
		return "", ErrKeyNotFound
	}

	s.log.Info("GET success", zap.String("key", key), zap.String("value", val))
	return val, nil
}

func (s *KVStore) Set(key, value string) error {
	s.log.Info("SET request", zap.String("key", key), zap.String("value", value))

	_, err := s.Conn.Do(
		tarantool.NewUpsertRequest("kvstore").
			Tuple([]interface{}{key, value}).
			Operations(tarantool.NewOperations().Assign(1, value)),
	).Get()

	if err != nil {
		s.log.Error("SET failed", zap.String("key", key), zap.Error(err))
	} else {
		s.log.Info("SET success", zap.String("key", key))
	}

	return err
}

func (s *KVStore) Delete(key string) error {
	s.log.Info("DELETE request", zap.String("key", key))

	data, err := s.Conn.Do(
		tarantool.NewDeleteRequest("kvstore").Key([]interface{}{key}),
	).Get()

	if err != nil {
		s.log.Error("DELETE failed", zap.String("key", key), zap.Error(err))
		return ErrInternal
	}

	if len(data) == 0 {
		s.log.Warn("DELETE key not found", zap.String("key", key))
		return ErrKeyNotFound
	}

	s.log.Info("DELETE success", zap.String("key", key))
	return nil
}
