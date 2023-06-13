package main

import (
	"fmt"
	"reverse/server"
	"rrbackend/local"
	"rrserver"
)

type DefaultHostURLMapper struct {
	DefaultHost string
}

func (mapper *DefaultHostURLMapper) MapURL(host string, url string) string {
	return "http://" + mapper.DefaultHost + url
}

func main() {
	fmt.Println("Request Response HTTP Proxy - Receiver")

	backend := local.RequestResponseBackend{}

	partiallyConfiguredServer := rrserver.SimpleRequestResponseServer{
		RequestChannelID: "request",
		Backend:          &backend,
	}

	reverseProxyServer := server.ReverseProxyServer{
		RRServer: partiallyConfiguredServer,
		UrlMapper: &DefaultHostURLMapper{
			DefaultHost: "localhost:3000",
		},
	}

	reverseProxyServer.Start()

}
