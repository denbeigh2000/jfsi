package jfsi

import (
	"github.com/satori/go.uuid"
)

type ID string

func NewID() ID {
	return ID(uuid.NewV4().String())
}
