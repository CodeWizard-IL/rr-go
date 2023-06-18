package main

import (
	"fmt"
	"forward/server"
	"rrbuilder"
	"time"
)

func main() {
	fmt.Println("Request Response HTTP Proxy - Forwarder")

	// Use ClientFromConfig to create a client from a config file

	clientConfig := rrbuilder.ClientConfig{
		Backend: rrbuilder.BackendConfig{
			Type: "amqp09",
			Configuration: map[string]any{
				"ConnectionString": "amqp://guest:guest@localhost:5672/",
			},
		},
		Type: "simple",
		Configuration: map[string]any{
			"TimeoutMillis":    10000,
			"RequestChannelID": "myrequest",
		},
	}

	client, err := rrbuilder.ClientFromConfig(clientConfig)
	if err != nil {
		fmt.Println("Error creating client: ", err)
		return
	}

	proxyServer := server.ForwardProxyServer{
		RRClient:       client,
		ListenAddress:  ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	proxyServer.Start()
}
