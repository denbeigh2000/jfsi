package application

import (
	"io"
	"log"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/application/chunker"
	"github.com/denbeigh2000/jfsi/metastore"
	"github.com/denbeigh2000/jfsi/storage"

	"github.com/satori/go.uuid"
)

type Node interface {
	Create(io.Reader) (jfsi.ID, error)
	Retrieve(jfsi.ID) (io.Reader, error)
	Update(jfsi.ID, io.Reader) error
	Delete(jfsi.ID) error
}

func NewNode(sc StorageConfig, c chunker.Chunker, ms metastore.MetaStore) Node {
	return node{StorageConfig: sc, Chunker: c, MetaStore: ms}
}

type node struct {
	StorageConfig StorageConfig
	Chunker       chunker.Chunker
	MetaStore     metastore.MetaStore
}

func (n node) key() jfsi.ID {
	return jfsi.ID(uuid.NewV4().String())
}

func (n node) createChunk(chunkID jfsi.ID, r io.Reader) error {
	log.Printf("Uploading chunk %v", chunkID)
	nodes := n.StorageConfig.Select(chunkID)

	err := storage.ParallelCreate(nodes, chunkID, r)
	if err != nil {
		return err
	}

	return nil
}

func (n node) Create(r io.Reader) (jfsi.ID, error) {
	id := n.key()
	log.Printf("Creating chunks for %v", id)
	chunks, err := n.Chunker.Chunk(r)
	if err != nil {
		return jfsi.ID(""), err
	}

	log.Printf("Creating metastore entries for %v", id)
	record, err := n.MetaStore.Create(id, len(chunks))
	if err != nil {
		return jfsi.ID(""), err
	}

	// TODO: parallelise this
	log.Printf("Uploading chunks for %v", id)
	for i, chunk := range record.Chunks {
		err = n.createChunk(chunk, chunks[i])
		if err != nil {
			return jfsi.ID(""), err
		}
	}

	return id, nil
}

func (n node) retrieveChunk(chunkID jfsi.ID) (io.Reader, error) {
	log.Printf("Retrieving chunk %v", chunkID)
	nodes := n.StorageConfig.Select(chunkID)
	r, err := storage.SelectiveRetrieve(nodes, chunkID)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (n node) Retrieve(id jfsi.ID) (io.Reader, error) {
	log.Printf("Retrieving metastore records for %v", id)
	record, err := n.MetaStore.Retrieve(id)
	if err != nil {
		return nil, err
	}

	chunkReaders := make([]io.Reader, len(record.Chunks))
	log.Printf("Retrieving chunks for %v", id)
	for i, chunk := range record.Chunks {
		r, err := n.retrieveChunk(chunk)
		if err != nil {
			return nil, err
		}

		chunkReaders[i] = r
	}

	return io.MultiReader(chunkReaders...), nil
}

func (n node) updateChunk(chunkID jfsi.ID, r io.Reader) error {
	log.Printf("Updating chunk %v", chunkID)
	nodes := n.StorageConfig.Select(chunkID)

	err := storage.ParallelUpdate(nodes, chunkID, r)
	if err != nil {
		return err
	}

	return nil
}

func (n node) Update(id jfsi.ID, r io.Reader) error {
	record, err := n.MetaStore.Retrieve(id)
	if err != nil {
		return err
	}

	chunks, err := n.Chunker.Chunk(r)
	if err != nil {
		return err
	}

	for _, chunkID := range record.Chunks {
		err = n.deleteChunk(chunkID)
		if err != nil {
			log.Printf("Error deleting chunk %v - proceeding anyway: %v", chunkID, err)
			err = nil
		}
	}

	log.Printf("Deleting old metastore record %v\n", id)
	err = n.MetaStore.Delete(id)
	if err != nil {
		return err
	}

	log.Printf("Creating new metastore record %v\n", id)
	record, err = n.MetaStore.Create(id, len(chunks))
	if err != nil {
		return err
	}

	log.Printf("Updating record %v, creating new chunks\n", id)
	for i, chunkID := range record.Chunks {
		err = n.createChunk(chunkID, chunks[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (n node) deleteChunk(id jfsi.ID) error {
	log.Printf("Deleting chunk %v\n", id)
	nodes := n.StorageConfig.Select(id)

	err := storage.ParallelDelete(nodes, id)
	if err != nil {
		return err
	}

	return nil
}

func (n node) Delete(id jfsi.ID) error {
	record, err := n.MetaStore.Retrieve(id)
	if err != nil {
		return err
	}

	log.Printf("Deleting chunks for %v", id)
	for _, chunkID := range record.Chunks {
		err = n.deleteChunk(chunkID)
		if err != nil {
			return err
		}
	}

	log.Printf("Deleting metastore record for %v", id)
	err = n.MetaStore.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
