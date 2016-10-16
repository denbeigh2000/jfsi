package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/denbeigh2000/jfsi/storage/disk"
	"github.com/denbeigh2000/jfsi/storage/handler"
)

var (
	port = flag.Int("port", 8080, "Port to serve on")
	dir  = flag.String("dir", "./.jfsi", "Directory to store blobs in")
)

func init() {
	flag.Parse()
}

func main() {
	store := disk.NewDiskStore(*dir)
	handler := handler.NewHTTP(store)

	host := fmt.Sprintf(":%v", *port)
	log.Printf("Storage serving on %v...", host)
	log.Fatal(http.ListenAndServe(host, handler))
}
