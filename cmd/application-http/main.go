package main

import (
	"net/http"

	"github.com/denbeigh2000/jfsi/application"
	"github.com/denbeigh2000/jfsi/application/handler"
	"github.com/denbeigh2000/jfsi/storage"
	"github.com/denbeigh2000/jfsi/storage/client"
)

func main() {
	store := client.NewClient("localhost", 8080)
	config := application.NewStorageConfig([]storage.Store{store})
	node := application.NewNode(config)
	handler := handler.NewHTTP(node)

	http.ListenAndServe(":8079", handler)
}
