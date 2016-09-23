package main

import (
	"net/http"

	"github.com/denbeigh2000/jfsi/storage/disk"
	"github.com/denbeigh2000/jfsi/storage/handler"
)

func main() {
	store := disk.NewDiskStore("./.jfsi")
	handler := handler.NewHTTP(store)

	http.ListenAndServe(":8080", handler)
}
