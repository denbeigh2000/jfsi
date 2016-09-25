package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/denbeigh2000/jfsi/application"
	"github.com/denbeigh2000/jfsi/application/handler"
	"github.com/denbeigh2000/jfsi/storage"
	"github.com/denbeigh2000/jfsi/storage/client"
)

var port = flag.Int("port", 8100, "Port to serve on")

func init() {
	flag.Parse()
}

func main() {
	stores := []storage.Store{
		client.NewClient("localhost", 8000),
		client.NewClient("localhost", 8001),
		client.NewClient("localhost", 8002),
		client.NewClient("localhost", 8003),
	}
	config := application.NewStorageConfig(stores)
	node := application.NewNode(config)
	handler := handler.NewHTTP(node)

	http.ListenAndServe(fmt.Sprintf(":%v", *port), handler)
}
