package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/denbeigh2000/jfsi/storage/cassandra"
	"github.com/denbeigh2000/jfsi/storage/handler"
)

var (
	port     = flag.Int("port", 8080, "Port to serve on")
	keyspace = flag.String("keyspace", "jfsi", "Keyspace to use")
	hostFlag arrayFlags
)

func init() {
	flag.Var(&hostFlag, "hosts", "Cassandra hosts")
	flag.Parse()
}

func main() {
	hosts := []string(hostFlag)
	store := cassandra.New(*keyspace, hosts...)
	handler := handler.NewHTTP(store)

	host := fmt.Sprintf(":%v", *port)
	log.Printf("Storage serving on %v...", host)
	log.Fatal(http.ListenAndServe(host, handler))
}
