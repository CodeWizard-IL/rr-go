package main

import (
	"fmt"
	"forward/server"
	"rrbackend/local"
	"rrclient"
	"time"
)

func main() {
	fmt.Println("Request Response HTTP Proxy - Forwarder")

	//TODO: make this configurable

	backend := local.RequestResponseBackend{}

	client := rrclient.SimpleRequestResponseClient{
		Backend:          &backend,
		TimeoutMillis:    10000,
		RequestChannelID: "request",
	}

	proxyServer := server.ForwardProxyServer{
		RRClient:       &client,
		ListenAddress:  ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	proxyServer.Start()
}
