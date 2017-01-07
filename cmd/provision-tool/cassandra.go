package main

import (
	"fmt"
	"strings"

	"github.com/denbeigh2000/jfsi/storage/cassandra"

	"github.com/gocql/gocql"
)

type CassandraProvisioner struct {
	Hosts       []string
	Keyspace    string
	Replication int
}

func (p CassandraProvisioner) Provision() error {
	query := fmt.Sprintf(cassandra.CreateKeyspace, p.Keyspace, p.Replication)
	cluster := gocql.NewCluster(p.Hosts...)
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	cluster.Keyspace = p.Keyspace

	err = session.Query(query).Exec()
	if err != nil {
		session.Close()
		return err
	}

	session.Close()

	session, err = cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Query(cassandra.CreateStorage).Exec()
	if err != nil {
		return err
	}

	return session.Query(cassandra.CreateMetastore).Exec()
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
