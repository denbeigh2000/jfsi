package disk

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/storage"
)

type store struct {
	dir string
}

func (s store) path(id jfsi.ID) string {
	if s.dir == "" {
		panic("disk store dir is nil - cannot continue")
	}

	return filepath.Join(s.dir, id.String())
}

func (s store) exists(id jfsi.ID) bool {
	path := s.path(id)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (s store) persist(id jfsi.ID, r io.Reader) error {
	path := s.path(id)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}

func (s store) read(id jfsi.ID) (io.Reader, error) {
	path := s.path(id)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)
	return r, nil
}

func (s store) Create(id jfsi.ID, r io.Reader) error {
	if s.exists(id) {
		return storage.AlreadyExistsErr(id)
	}

	err := s.persist(id, r)
	if err != nil {
		return err
	}

	return nil
}

func (s store) Retrieve(id jfsi.ID) (io.Reader, error) {
	if !s.exists(id) {
		return nil, storage.NotFoundErr(id)
	}

	r, err := s.read(id)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s store) Update(id jfsi.ID, r io.Reader) error {
	if !s.exists(id) {
		return storage.NotFoundErr(id)
	}

	err := s.persist(id, r)
	if err != nil {
		return err
	}

	return nil
}

func (s store) Delete(id jfsi.ID) error {
	if !s.exists(id) {
		return storage.NotFoundErr(id)
	}

	path := s.path(id)
	return os.Remove(path)
}

func NewDiskStore(dir string) storage.Store {
	// err := os.Mkdir(dir, os.ModePerm)
	// if err != nil && !os.IsExist(err) {
	// 	panic(err)
	// }

	return store{
		dir: dir,
	}
}
