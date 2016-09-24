package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func logRequest(req *http.Request) {
	dump, err := httputil.DumpRequest(req, false)
	if err != nil {
		log.Printf("Dumping request failed: %v", err)
	}

	log.Printf("%v", string(dump))
}
