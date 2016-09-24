package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/denbeigh2000/jfsi/storage/disk"
	"github.com/denbeigh2000/jfsi/storage/handler"
)

var (
	port = flag.Int("port", 8080, "Port to serve on")
	dir  = flag.String("port", "./.jfsi", "Directory to store blobs in")
)

func init() {
	flag.Parse()
}

func main() {
	store := disk.NewDiskStore(*dir)
	handler := handler.NewHTTP(store)

	http.ListenAndServe(fmt.Sprintf(":%v", *port), handler)
}
