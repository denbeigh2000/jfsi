package metastore

import (
	"github.com/denbeigh2000/jfsi"
)

type Record struct {
	Key    jfsi.ID
	Chunks []jfsi.ID
}

type MetaStore interface {
	Create(key jfsi.ID, n int) (Record, error)
	Retrieve(key jfsi.ID) (Record, error)
	Update(key jfsi.ID, r Record) error
	Delete(key jfsi.ID) error
}
