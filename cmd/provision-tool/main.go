package main

import (
	"flag"
	"log"
)

var (
	keyspace    = flag.String("keyspace", "jfsi", "Keyspace to provision")
	replication = flag.Int("replication", 2, "Replication factor to use")
	hostFlag    arrayFlags
)

func init() {
	flag.Var(&hostFlag, "hosts", "Cassandra hosts")
	flag.Parse()
}

func main() {
	hosts := []string(hostFlag)
	provisioner := CassandraProvisioner{
		Hosts:       hosts,
		Keyspace:    *keyspace,
		Replication: *replication,
	}

	err := provisioner.Provision()
	if err != nil {
		log.Fatalf("Provisioning failed: %v", err)
	}

	log.Println("Successfully provisioned!")
}
