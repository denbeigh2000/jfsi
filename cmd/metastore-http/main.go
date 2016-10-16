package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/denbeigh2000/jfsi/metastore/handler"
	"github.com/denbeigh2000/jfsi/metastore/memory"
)

var port = flag.Int("port", 8200, "Port to serve on")

func init() {
	flag.Parse()
}

func main() {
	store := memory.NewStore()
	handler := handler.NewHTTP(store)

	host := fmt.Sprintf(":%v", *port)

	fmt.Printf("Metastore serving on %v...\n", host)
	log.Fatal(http.ListenAndServe(host, handler))
}
