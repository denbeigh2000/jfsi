package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/denbeigh2000/jfsi/metastore/handler"
	"github.com/denbeigh2000/jfsi/metastore/redis"
)

var (
	port          = flag.Int("port", 8200, "Port to serve on")
	redisPort     = flag.Int("redis-port", 6379, "Redis backend port to use")
	redisHost     = flag.String("redis-host", "localhost", "Redis host port to use")
	redisDB       = flag.Int("redis-db", 0, "Redis DB to use")
	redisPassword = flag.String("redis-pw", "", "Redis password to use")
)

func init() {
	flag.Parse()
}

func main() {
	redisAddr := fmt.Sprintf("%v:%v", *redisHost, *redisPort)
	store := redis.NewStore(redisAddr, *redisPassword, 0)
	handler := handler.NewHTTP(store)

	host := fmt.Sprintf(":%v", *port)

	fmt.Printf("Metastore serving on %v...\n", host)
	log.Fatal(http.ListenAndServe(host, handler))
}
