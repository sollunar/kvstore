package kvstore

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

type KVService struct {
	store Storage
}

func NewKVService(s Storage) *KVService {
	return &KVService{store: s}
}

func (s *KVService) Get(key string) (string, error) {
	return s.store.Get(key)
}

func (s *KVService) Set(key, value string) error {
	return s.store.Set(key, value)
}

func (s *KVService) Delete(key string) error {
	return s.store.Delete(key)
}
