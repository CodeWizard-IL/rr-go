package main

import (
	fserver "forward/server"
	"log"
	rserver "reverse/server"
	//"rrbackend/local"
	. "rrbackendazsmb"
	"rrclient"
	"rrserver"
	"time"
)

type DefaultHostURLMapper struct {
	DefaultHost string
}

func (mapper *DefaultHostURLMapper) MapURL(_ string, url string) string {
	return "http://" + mapper.DefaultHost + url
}

func main() {
	//backend := local.RequestResponseBackend{}

	backend := RRBackendAzSMB{
		ConnectionString:  "Endpoint=sb://cwalexeyrr.servicebus.windows.net/;SharedAccessKeyName=rrgo;SharedAccessKey=sKMyUVlVxhjG62QrJh3mLlS/zXLpIK/a9+ASbLD88Xc=",
		RequestQueueName:  "myrequest",
		ResponseQueueName: "myrequest-response",
	}

	err := backend.Connect()
	if err != nil {
		log.Fatal("Error connecting to backend: ", err)
	}

	client := rrclient.SimpleRequestResponseClient{
		Backend:          &backend,
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
