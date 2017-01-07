package metastore

import (
	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/metastore"

	"github.com/gocql/gocql"
)

func NewStore(keyspace string, hosts ...string) metastore.MetaStore {
	store := &metaStore{
		Hosts:    hosts,
		Keyspace: keyspace,
	}
	store.configure()

	return store
}

type metaStore struct {
	Hosts    []string
	Keyspace string

	cluster *gocql.ClusterConfig
}

func (s *metaStore) configure() {
	s.cluster = gocql.NewCluster(s.Hosts...)
	s.cluster.Keyspace = s.Keyspace
}

func (s metaStore) session() (*gocql.Session, error) {
	session, err := s.cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	session.SetConsistency(gocql.LocalOne)
	return session, nil
}

var (
	createQuery   = `INSERT INTO chunk (key, chunks) VALUES (?, ?);`
	retrieveQuery = `SELECT chunks FROM chunk WHERE key = ?;`
	updateQuery   = `UPDATE chunk SET chunks = ? FROM metastore WHERE key = ?;`
	deleteQuery   = `DELETE FROM chunk WHERE key = ?;`
)

func (m metaStore) Create(key jfsi.ID, n int) (r metastore.Record, err error) {
	session, err := m.session()
	if err != nil {
		return
	}
	defer session.Close()

	uuids, cqluuids := make([]jfsi.ID, n, n), make([]gocql.UUID, n, n)
	for i := 0; i < n; i++ {
		uuids[i] = jfsi.NewID()
		cqluuids[i] = gocql.UUID(uuids[i])
	}

	err = session.Query(createQuery, gocql.UUID(key), cqluuids).Exec()
	if err != nil {
		return
	}

	r.Chunks = uuids
	r.Key = key

	return
}

func (m metaStore) Retrieve(key jfsi.ID) (r metastore.Record, err error) {
	session, err := m.session()
	if err != nil {
		return
	}
	defer session.Close()

	chunks := make([]gocql.UUID, 0)
	err = session.Query(retrieveQuery, gocql.UUID(key)).Scan(&chunks)
	if err != nil {
		return
	}

	nativeChunks := make([]jfsi.ID, len(chunks))

	for i, chunk := range chunks {
		nativeChunks[i] = jfsi.ID(chunk)
	}

	r.Chunks = nativeChunks
	r.Key = key

	return
}

func (m metaStore) Update(key jfsi.ID, r metastore.Record) error {
	session, err := m.session()
	if err != nil {
		return err
	}
	defer session.Close()

	chunks := make([]gocql.UUID, len(r.Chunks))
	for i, chunk := range r.Chunks {
		chunks[i] = gocql.UUID(chunk)
	}

	if err = session.Query(updateQuery, gocql.UUID(key), chunks).Exec(); err != nil {
		return err
	}

	return nil
}

func (m metaStore) Delete(key jfsi.ID) error {
	session, err := m.session()
	if err != nil {
		return err
	}
	defer session.Close()

	if err = session.Query(deleteQuery, gocql.UUID(key)).Exec(); err != nil {
		return err
	}

	return nil
}
