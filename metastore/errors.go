package metastore

import (
	"fmt"

	"github.com/denbeigh2000/jfsi"
)

const ZeroLenStr = "Cannot allocate MetaStore entry with zero-length capacity"

type KeyAlreadyExistsErr jfsi.ID

func (err KeyAlreadyExistsErr) Error() string {
	return fmt.Sprintf("Key already exists: %v", jfsi.ID(err))
}

type KeyNotFoundErr jfsi.ID

func (err KeyNotFoundErr) Error() string {
	return fmt.Sprintf("Key not found: %v", jfsi.ID(err))
}

type ZeroLengthCapacityRecordErr struct{}

func (err ZeroLengthCapacityRecordErr) Error() string {
	return ZeroLenStr
}
