package application

import (
	"hash/crc64"
	"io"
	"log"
	"math/rand"
	"sync"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/storage"
)

var table = crc64.MakeTable(crc64.ISO)

func Select(stores []storage.Store) storage.Store {
	n := rand.Int31n(int32(len(stores)))
	return stores[n]
}

// TODO: Move these to storage layer
func NewDiskStorageConfig(stores []storage.Store, replication int) DiskStorageConfig {
	return DiskStorageConfig{
		RWMutex:     &sync.RWMutex{},
		Replication: replication,
		Stores:      stores,
	}
}

type DiskStorageConfig struct {
	*sync.RWMutex

	Replication int
	Stores      []storage.Store
}

type Selecter interface {
	Select(jfsi.ID) []storage.Store
}

type Transferrer interface {
	Transfer(jfsi.ID, []storage.Store, io.Reader) error
}

func (s *DiskStorageConfig) Select(id jfsi.ID) []storage.Store {
	s.RLock()
	defer s.RUnlock()
	n := len(s.Stores)
	if s.Replication >= n {
		log.Panicf("replication factor (%v) must be less than number of storage nodes (%v)", s.Replication, n)
	}

	hash := crc64.Checksum([]byte(id.String()), table)
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
