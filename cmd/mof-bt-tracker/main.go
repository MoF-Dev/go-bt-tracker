package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("failed to start http server: %s", err)
	}
}
