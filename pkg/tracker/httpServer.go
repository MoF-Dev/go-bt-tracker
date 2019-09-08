package tracker

import (
	"github.com/MoF-Dev/go-bt-tracker/internal/app/http"
	"github.com/gorilla/mux"
	"net"
	http2 "net/http"
)

type httpServerExtras interface {
}

type HttpServer interface {
	Server
	httpServerExtras
}

func ListenHttp(server HttpServer, listener net.Listener) error {
	r := mux.NewRouter()
	r.HandleFunc("/announce", http.AnnounceHandler(server))
	r.HandleFunc("/scrape", http.ScrapeHandler(server))
	return http2.Serve(listener, r)
}
