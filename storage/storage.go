package storage

import (
	"fmt"
	"io"

	"github.com/denbeigh2000/jfsi"
)

type AlreadyExistsErr jfsi.ID

func (err AlreadyExistsErr) Error() string {
	return fmt.Sprintf("Key already exists: %v", string(err))
}

type NotFoundErr jfsi.ID

func (err NotFoundErr) Error() string {
	return fmt.Sprintf("Key not found: %v", string(err))
}

type Readerer interface {
	Reader() io.Reader
}

type Item struct {
	Readerer
	ID jfsi.ID
}

type Store interface {
	Create(jfsi.ID, io.Reader) error
	Retrieve(jfsi.ID) (io.Reader, error)
	Update(jfsi.ID, io.Reader) error
	Delete(jfsi.ID) error
}
