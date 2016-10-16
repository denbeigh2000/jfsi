package metastore

import (
	"fmt"

	"github.com/denbeigh2000/jfsi"
)

type KeyAlreadyExistsErr jfsi.ID

func (err KeyAlreadyExistsErr) Error() string {
	return fmt.Sprintf("Key already exists: %v", string(err))
}

type KeyNotFoundErr jfsi.ID

func (err KeyNotFoundErr) Error() string {
	return fmt.Sprintf("Key not found: %v", string(err))
}

type ZeroLengthCapacityRecordErr struct{}

func (err ZeroLengthCapacityRecordErr) Error() string {
	return "Cannot allocate MetaStore entry with zero-length capacity"
}
