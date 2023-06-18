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
	err := server.RRClient.Start()
	if err != nil {
		log.Fatal("Error starting client: ", err)
	}

	proxyHandler := handler.NewProxyHandler(server.RRClient)

	s := &http.Server{
		Addr:           server.ListenAddress,
		Handler:        proxyHandler,
		ReadTimeout:    server.ReadTimeout,
		WriteTimeout:   server.WriteTimeout,
		MaxHeaderBytes: server.MaxHeaderBytes,
	}
	log.Fatal(s.ListenAndServe())
}
