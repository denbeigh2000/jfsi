package cassandra

var CreateStorage = `
	CREATE TABLE storage (
		key uuid PRIMARY KEY,
		data blob
	);
`

// Denormalised storage to persist chunk mappings
var CreateMetastore = `
	CREATE TABLE chunk (
		key uuid PRIMARY KEY,
		chunks list<uuid>
	);
`

var CreateKeyspace = `
	CREATE KEYSPACE %v WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': %d};
`
