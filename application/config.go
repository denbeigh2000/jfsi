package application

import (
	"encoding/binary"
	"log"
	"sync"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/storage"
)

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

	hash, _ := binary.Uvarint([]byte(id))
	mod := int(hash) % n

	stores := make([]storage.Store, 1+s.Replication)
	for i := 0; i <= s.Replication; i++ {
		stores[i] = s.Stores[mod+i]
	}

	return stores
}
