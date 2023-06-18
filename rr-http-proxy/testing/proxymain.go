package main

import (
	fserver "forward/server"
	"log"
	"reverse/processor"
	rserver "reverse/server"
	"rrbuilder"
	"rrclient"
	"rrserver"
	"time"
)

func main() {
	//backend := local.RequestResponseBackend{}

	backendConfig := rrbuilder.BackendConfig{
		Type: "azsb",
		Configuration: map[string]any{
			"ConnectionString":    "Endpoint=sb://cwalexeyrr.servicebus.windows.net/;SharedAccessKeyName=rrgo;SharedAccessKey=sKMyUVlVxhjG62QrJh3mLlS/zXLpIK/a9+ASbLD88Xc=",
			"RequestQueueName":    "myrequest",
			"ResponseQueueName":   "myrequest-response",
			"MinSessionReceivers": 1,
		},
	}

	//backendConfig := rrbuilder.BackendConfig{
	//	Type:          "local",
	//	Configuration: map[string]any{},
	//}

	//backendConfig := rrbuilder.BackendConfig{
	//	Type: "amqp09",
	//	Configuration: map[string]any{
	//		"ConnectionString": "amqp://guest:guest@localhost:5672/",
	//	},
	//}

	backend, err := rrbuilder.BackendFromConfig(backendConfig)
	if err != nil {
		log.Fatal("Error creating backend: ", err)
	}

	err = backend.Connect()
	if err != nil {
		log.Fatal("Error connecting to backend: ", err)
	}

	client := rrclient.SimpleRequestResponseClient{
		Backend:          backend,
		TimeoutMillis:    10000,
		RequestChannelID: "myrequest",
	}

	proxyServer := fserver.ForwardProxyServer{
		RRClient:       &client,
		ListenAddress:  ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	partiallyConfiguredServer := rrserver.SimpleRequestResponseServer{
		RequestChannelID: "myrequest",
		Backend:          backend,
	}

	reverseProxyServer := rserver.ReverseProxyServer{
		RRServer: partiallyConfiguredServer,
		//UrlMapper: &DefaultHostURLMapper{
		//	DefaultHost: "localhost:3000",
		//},
		UrlMapper: &processor.FirstPathUrlMapper{},
	}

	reverseProxyServer.Start()

	proxyServer.Start()
}
