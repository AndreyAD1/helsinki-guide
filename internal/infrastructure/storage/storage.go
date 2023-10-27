package storage

type Repository interface {
	GetAllAdresses() ([]string, error)
}

type MemoryStorage struct {
	store map[string]string
}

func (ms MemoryStorage) GetAllAdresses() ([]string, error) {
	addresses := []string{}
	for address := range ms.store {
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func NewStorage(dbURL string) (Repository, error) {
	return MemoryStorage{make(map[string]string)}, nil
}
