package server

import (
	"reverse/processor"
	"rr-lib/rrserver"
)

type ReverseProxyServer struct {
	RRServer  rrserver.SimpleRequestResponseServer
	UrlMapper processor.URLMapper
}

func (server *ReverseProxyServer) Start() {
	server.RRServer.Processor = &processor.ReverseProxyProcessor{UrlMapper: server.UrlMapper}
	server.RRServer.Start()
}
