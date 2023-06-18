package main

import (
	"fmt"
	"reverse/server"
	"rrbuilder"
	"rrserver"
)

type DefaultHostURLMapper struct {
	DefaultHost string
}

func (mapper *DefaultHostURLMapper) MapURL(_ string, url string) string {
	return "http://" + mapper.DefaultHost + url
}

func main() {
	fmt.Println("Request Response HTTP Proxy - Receiver")

	// Use ServerFromConfig to create a server from a config file

	serverConfig := rrbuilder.ServerConfig{
		Backend: rrbuilder.BackendConfig{
			Type: "amqp09",
			Configuration: map[string]any{
				"ConnectionString": "amqp://guest:guest@localhost:5672/",
			},
		},
		Type: "simple",
		Configuration: map[string]any{
			"RequestChannelID": "myrequest",
		},
	}

	partiallyConfiguredServer, err := rrbuilder.ServerFromConfig(serverConfig)
	if err != nil {
		fmt.Println("Error creating server: ", err)
		return
	}

	// Assert that the server is a SimpleRequestResponseServer

	simpleServer, ok := partiallyConfiguredServer.(*rrserver.SimpleRequestResponseServer)
	if !ok {
		fmt.Println("Error: Server is not a SimpleRequestResponseServer")
		return
	}

	reverseProxyServer := server.ReverseProxyServer{
		RRServer: *simpleServer,
		UrlMapper: &DefaultHostURLMapper{
			DefaultHost: "localhost:3000",
		},
	}

	reverseProxyServer.Start()

	// Wait forever
	select {}

}
