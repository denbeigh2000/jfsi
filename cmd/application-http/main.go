package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/denbeigh2000/jfsi/application"
	"github.com/denbeigh2000/jfsi/application/chunker"
	"github.com/denbeigh2000/jfsi/application/handler"
	msClient "github.com/denbeigh2000/jfsi/metastore/client"
	"github.com/denbeigh2000/jfsi/storage"
	"github.com/denbeigh2000/jfsi/storage/client"
)

var port = flag.Int("port", 8100, "Port to serve on")
var replication = flag.Int("replication", 1, "Replication factor to store chunks with")
var chunkSize = flag.Int("chunk-size", 131072, "Chunk size in bytes (default 128KiB)")

func init() {
	flag.Parse()
}

func main() {
	stores := []storage.Store{
		client.NewClient("localhost", 8000),
		//client.NewClient("localhost", 8001),
		//client.NewClient("localhost", 8002),
		//client.NewClient("localhost", 8003),
	}
	// config := application.NewStorageConfig(stores, *replication)
	chunker := chunker.NewChunker(int64(*chunkSize))
	metastore := msClient.NewHTTP("localhost", 8200)
	node := application.NewNode(stores, chunker, metastore)
	handler := handler.NewHTTP(node)

	host := fmt.Sprintf(":%v", *port)
	log.Printf("Application serving on %v...", host)
	log.Fatal(http.ListenAndServe(host, handler))
}
