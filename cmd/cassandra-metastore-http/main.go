package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/denbeigh2000/jfsi/cmd/common"
	"github.com/denbeigh2000/jfsi/metastore/cassandra"
	"github.com/denbeigh2000/jfsi/metastore/handler"
)

var (
	port     = flag.Int("port", 8200, "Port to serve on")
	keyspace = flag.String("keyspace", "jfsi", "Cassandra keyspace to use")
	hostFlag common.ArrayFlags
)

func init() {
	flag.Var(&hostFlag, "hosts", "Cassandra hosts")
	flag.Parse()
}

func main() {
	store := metastore.NewStore(*keyspace, hostFlag...)
	handler := handler.NewHTTP(store)

	host := fmt.Sprintf(":%v", *port)

	fmt.Printf("Metastore serving on %v...\n", host)
	log.Fatal(http.ListenAndServe(host, handler))
}
