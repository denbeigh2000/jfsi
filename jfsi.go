package jfsi

import (
	"github.com/satori/go.uuid"
)

type ID [16]byte

func NewID() ID {
	return ID(uuid.NewV4())
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}

func IDFromString(id string) (ID, error) {
	u, err := uuid.FromString(id)
	if err != nil {
		return ID{}, err
	}

	return ID(u), nil
}
