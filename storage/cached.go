package storage

import (
	"io"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/cache"
)

type CachedStore struct {
	Store Store
	Cache cache.Cache
}

func parallelWrite(ID jfsi.ID, r io.Reader, c cache.Cache, s Store) error {
	pr, pw := io.Pipe()
	tr := io.TeeReader(r, pw)

	go func() {
		defer pw.Close()
		c.Set(ID, pr)
	}()

	return s.Create(ID, tr)
}

func (c CachedStore) Create(ID jfsi.ID, r io.Reader) error {
	return parallelWrite(ID, r, c.Cache, c.Store)
}

func (c CachedStore) Update(ID jfsi.ID, r io.Reader) error {
	return parallelWrite(ID, r, c.Cache, c.Store)
}

func (c CachedStore) Delete(ID jfsi.ID) error {
	c.Cache.Delete(ID)
	return c.Store.Delete(ID)
}

func (c CachedStore) Retrieve(ID jfsi.ID) (io.Reader, error) {
	b, ok := c.Cache.Get(ID)
	if ok {
		return b, nil
	}

	return c.Store.Retrieve(ID)
}
