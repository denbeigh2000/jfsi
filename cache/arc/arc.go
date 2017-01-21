package arc

import (
	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/cache"

	"github.com/hashicorp/golang-lru"

	"bytes"
	"io"
	"io/ioutil"
)

func NewARCCache(size int) (cache.Cache, error) {
	impl, err := lru.NewARC(size)
	if err != nil {
		return nil, err
	}

	return &ARC{impl: impl}, nil
}

type ARC struct {
	impl *lru.ARCCache
}

func (a *ARC) Get(key jfsi.ID) (io.Reader, bool) {
	i, ok := a.impl.Get(key)
	if !ok {
		return nil, false
	}

	b, ok := i.([]byte)
	if !ok {
		return nil, false
	}

	return bytes.NewReader(b), true
}

func (a *ARC) Set(key jfsi.ID, r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	a.impl.Add(key, b)
	return nil
}
