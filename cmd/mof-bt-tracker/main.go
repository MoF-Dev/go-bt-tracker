package main

import (
	http2 "github.com/MoF-Dev/go-bt-tracker/internal/app/http"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/announce", http2.AnnounceHandler)
	r.HandleFunc("/scrape", http2.ScrapeHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("failed to start http server: %s", err)
	}
}
