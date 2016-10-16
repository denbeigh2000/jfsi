package memory

import (
	"sync"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/metastore"
)

func NewStore() metastore.MetaStore {
	return &store{
		items: make(map[jfsi.ID][]jfsi.ID),
	}
}

type store struct {
	sync.RWMutex

	items map[jfsi.ID][]jfsi.ID
}

func (s *store) Create(key jfsi.ID, n int) (r metastore.Record, err error) {
	if n <= 0 {
		err = metastore.ZeroLengthCapacityRecordErr{}
		return
	}
	s.Lock()
	defer s.Unlock()

	_, ok := s.items[key]
	if ok {
		err = metastore.KeyAlreadyExistsErr(key)
		return
	}

	items := make([]jfsi.ID, n)
	for i := 0; i < n; i++ {
		items[i] = jfsi.NewID()
	}

	s.items[key] = items

	r = metastore.Record{
		Key:    key,
		Chunks: items,
	}
	return
}

func (s *store) Retrieve(key jfsi.ID) (r metastore.Record, err error) {
	s.RLock()
	defer s.RUnlock()

	item, ok := s.items[key]
	if !ok {
		err = metastore.KeyNotFoundErr(key)
		return
	}

	r = metastore.Record{
		Key:    key,
		Chunks: item,
	}
	return
}

func (s *store) Update(key jfsi.ID, r metastore.Record) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.items[key]
	if !ok {
		return metastore.KeyNotFoundErr(key)
	}

	s.items[key] = r.Chunks
	return nil
}

func (s *store) Delete(key jfsi.ID) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.items[key]
	if !ok {
		return metastore.KeyNotFoundErr(key)
	}

	delete(s.items, key)
	return nil
}
