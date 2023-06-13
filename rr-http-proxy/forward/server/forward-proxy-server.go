package server

import (
	"forward/handler"
	"log"
	"net/http"
	"rrclient"
	"time"
)

type ForwardProxyServer struct {
	RRClient       rrclient.RequestResponseClient
	ListenAddress  string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

func (server *ForwardProxyServer) Start() {
	handler := handler.NewProxyHandler(server.RRClient)

	s := &http.Server{
		Addr:           server.ListenAddress,
		Handler:        handler,
		ReadTimeout:    server.ReadTimeout,
		WriteTimeout:   server.WriteTimeout,
		MaxHeaderBytes: server.MaxHeaderBytes,
	}
	log.Fatal(s.ListenAndServe())
}
