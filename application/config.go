package application

import (
	"hash/crc64"
	"log"
	"sync"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/storage"
)

var table = crc64.MakeTable(crc64.ISO)

func NewStorageConfig(stores []storage.Store, replication int) StorageConfig {
	return StorageConfig{
		RWMutex:     &sync.RWMutex{},
		Replication: replication,
		Stores:      stores,
	}
}

type StorageConfig struct {
	*sync.RWMutex

	Replication int
	Stores      []storage.Store
}

func (s *StorageConfig) Select(id jfsi.ID) []storage.Store {
	s.RLock()
	defer s.RUnlock()
	n := len(s.Stores)
	if s.Replication >= n {
		log.Panicf("replication factor (%v) must be less than number of storage nodes (%v)", s.Replication, n)
	}

	hash := crc64.Checksum([]byte(id), table)
	mod := int(hash % uint64(n))

	if mod < 0 {
		mod = mod * -1
	}

	stores := make([]storage.Store, 1+s.Replication)
	for i := 0; i <= s.Replication; i++ {
		target := (mod + i) % n
		stores[i] = s.Stores[target]
	}

	return stores
}
