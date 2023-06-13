package main

import (
	fserver "forward/server"
	rserver "reverse/server"
	"rrbackend/local"
	"rrclient"
	"rrserver"
	"time"
)

type DefaultHostURLMapper struct {
	DefaultHost string
}

func (mapper *DefaultHostURLMapper) MapURL(host string, url string) string {
	return "http://" + mapper.DefaultHost + url
}

func main() {
	backend := local.RequestResponseBackend{}

	backend.Connect()

	client := rrclient.SimpleRequestResponseClient{
		Backend:          &backend,
		TimeoutMillis:    10000,
		RequestChannelID: "request",
	}

	proxyServer := fserver.ForwardProxyServer{
		RRClient:       &client,
		ListenAddress:  ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	partiallyConfiguredServer := rrserver.SimpleRequestResponseServer{
		RequestChannelID: "request",
		Backend:          &backend,
	}

	reverseProxyServer := rserver.ReverseProxyServer{
		RRServer: partiallyConfiguredServer,
		UrlMapper: &DefaultHostURLMapper{
			DefaultHost: "localhost:3000",
		},
	}

	reverseProxyServer.Start()

	proxyServer.Start()
}
