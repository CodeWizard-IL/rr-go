package server

import (
	"github.com/CodeWizard-IL/rr-go/rr-lib/rrserver"
	"rr-http-proxy/reverse/processor"
)

type ReverseProxyServer struct {
	RRServer  rrserver.SimpleRequestResponseServer
	UrlMapper processor.URLMapper
}

func (server *ReverseProxyServer) Start() {
	server.RRServer.Processor = &processor.ReverseProxyProcessor{UrlMapper: server.UrlMapper}
	server.RRServer.Start()
}
