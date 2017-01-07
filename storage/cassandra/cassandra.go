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

func uuidFromID(id jfsi.ID) (gocql.UUID, error) {
	return gocql.UUIDFromBytes([]byte(string(id)))
}

type store struct {
	hosts       []string
	keyspace    string
	consistency Consistency

	cluster *gocql.ClusterConfig
}

func New(keyspace, hosts ...string) (storage.Store, error) {
	newStore := &store{
		hosts:    hosts,
		keyspace: keyspace,
	}

	clusterConfig := gocql.NewCluster(hosts...)
}

func (s *store) configure() {
	s.cluster = gocql.NewCluster(s.hosts...)
	s.cluster.Keyspace = s.keyspace
}

func (s *store) session() (*gocql.Session, error) {
	return s.cluster.CreateSession()
}

func (s *store) Create(id jfsi.ID, r io.Reader) error {
	uuid, err := uuidFromID(id)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	session, err := s.session()
	if err != nil {
		return err
	}

	if err := session.Query(insertQuery, uuid, data).Exec(); err != nil {
		return err
	}

	return nil
}

func (s *store) Retrieve(id jfsi.ID) (io.Reader, error) {
	uuid, err := uuidFromID(id)
	if err != nil {
		return nil, err
	}

	session, err := s.session()
	if err != nil {
		return nil, err
	}

	query := session.Query(retrieveQuery, uuid)

	blob := []byte{}
	err = query.Scan(&blob)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(blob), nil
}

func (s *store) Update(id jfsi.ID, r io.Reader) error {
	uuid, err := uuidFromID(id)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	session, err := s.session()
	if err != nil {
		return err
	}

	if err := session.Query(updateQuery, uuid, data).Exec(); err != nil {
		return err
	}

	return nil
}

func (s *store) Delete(id jfsi.ID) error {
	uuid, err := uuidFromID(id)
	if err != nil {
		return err
	}

	session, err := s.session()
	if err != nil {
		return err
	}

	if err := session.Query(deleteQuery, uuid).Exec(); err != nil {
		return err
	}

	return nil
}
