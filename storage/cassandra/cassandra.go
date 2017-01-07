package cassandra

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/storage"

	"github.com/gocql/gocql"
)

type Consistency uint16

const (
	Any         Consistency = 0x00
	One         Consistency = 0x01
	Two         Consistency = 0x02
	Three       Consistency = 0x03
	Quorum      Consistency = 0x04
	All         Consistency = 0x05
	LocalQuorum Consistency = 0x06
	EachQuorum  Consistency = 0x07
	LocalOne    Consistency = 0x0A
)

const (
	insertQuery   = `INSERT INTO storage (key, data) VALUES (?, ?);`
	retrieveQuery = `SELECT data FROM storage WHERE key = ?;`
	updateQuery   = `UPDATE storage SET data = ? WHERE key = ?;`
	deleteQuery   = `DELETE FROM storage WHERE key = ?;`
)

type store struct {
	hosts       []string
	keyspace    string
	consistency Consistency

	cluster *gocql.ClusterConfig
}

func New(keyspace string, hosts ...string) storage.Store {
	newStore := &store{
		hosts:    hosts,
		keyspace: keyspace,
	}

	newStore.configure()
	return newStore
}

func (s *store) configure() {
	s.cluster = gocql.NewCluster(s.hosts...)
	s.cluster.Keyspace = s.keyspace
}

func (s *store) session() (*gocql.Session, error) {
	// TODO: Allow for more flexibility
	session, err := s.cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	session.SetConsistency(gocql.LocalOne)

	return session, nil
}

func (s *store) Create(id jfsi.ID, r io.Reader) error {
	uuid := gocql.UUID(id)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	session, err := s.session()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Query(insertQuery, uuid, data).Exec()
}

func (s *store) Retrieve(id jfsi.ID) (io.Reader, error) {
	uuid := gocql.UUID(id)
	session, err := s.session()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	query := session.Query(retrieveQuery, uuid)

	blob := []byte{}
	err = query.Scan(&blob)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(blob), nil
}

func (s *store) Update(id jfsi.ID, r io.Reader) error {
	uuid := gocql.UUID(id)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	session, err := s.session()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Query(updateQuery, uuid, data).Exec()
}

func (s *store) Delete(id jfsi.ID) error {
	uuid := gocql.UUID(id)
	session, err := s.session()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Query(deleteQuery, uuid).Exec()
}
