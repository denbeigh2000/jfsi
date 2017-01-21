package cache

import (
	"github.com/denbeigh2000/jfsi"

	"io"
)

type Cache interface {
	Get(key jfsi.ID) (r io.Reader, ok bool)
	Set(key jfsi.ID, r io.Reader) error
	Delete(key jfsi.ID)
}
