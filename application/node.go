package application

import (
	"io"

	"github.com/denbeigh2000/jfsi"

	"github.com/satori/go.uuid"
)

type Node interface {
	Create(io.Reader) (jfsi.ID, error)
	Retrieve(jfsi.ID) (io.Reader, error)
	Update(jfsi.ID, io.Reader) error
	Delete(jfsi.ID) error
}

func NewNode(sc StorageConfig) Node {
	return node{StorageConfig: sc}
}

type node struct {
	StorageConfig StorageConfig
}

func (n node) key() jfsi.ID {
	return jfsi.ID(uuid.NewV4().String())
}

func (n node) Create(r io.Reader) (jfsi.ID, error) {
	id := n.key()
	nodes := n.StorageConfig.Select(id)
	if len(nodes) != 1 {
		panic("Replication factor of >1 not yet supported")
	}

	err := nodes[0].Create(id, r)
	if err != nil {
		return jfsi.ID(""), err
	}

	return id, nil
}

func (n node) Retrieve(id jfsi.ID) (io.Reader, error) {
	nodes := n.StorageConfig.Select(id)
	if len(nodes) != 1 {
		panic("Replication factor of >1 not yet supported")
	}

	r, err := nodes[0].Retrieve(id)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (n node) Update(id jfsi.ID, r io.Reader) error {
	nodes := n.StorageConfig.Select(id)
	if len(nodes) != 1 {
		panic("Replication factor of >1 not yet supported")
	}

	err := nodes[0].Update(id, r)
	if err != nil {
		return err
	}

	return nil
}

func (n node) Delete(id jfsi.ID) error {
	nodes := n.StorageConfig.Select(id)
	if len(nodes) != 1 {
		panic("Replication factor of >1 not yet supported")
	}

	err := nodes[0].Delete(id)
	if err != nil {
		return err
	}

	return nil
}
