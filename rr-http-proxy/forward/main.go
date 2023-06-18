package main

import (
	"common/util"
	"fmt"
	"forward/server"
	"rrbuilder"
	"time"
)

type ListenerConfig struct {
	ListenAddress  string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type ForwardProxyConfig struct {
	Client   rrbuilder.ClientConfig
	Listener ListenerConfig
}

func main() {
	fmt.Println("Request Response HTTP Proxy - Forwarder")

	// Read ForwardProxyConfig from YAML file
	var config ForwardProxyConfig

	defaultConfigYaml := "config.yaml"

	err := util.ReadYamlToStruct(defaultConfigYaml, &config)
	if err != nil {
		return
	}
	clientConfig := config.Client

	client, err := rrbuilder.ClientFromConfig(clientConfig)
	if err != nil {
		fmt.Println("Error creating client: ", err)
		return
	}

	proxyServer := server.ForwardProxyServer{
		RRClient:       client,
		ListenAddress:  config.Listener.ListenAddress,
		ReadTimeout:    config.Listener.ReadTimeout,
		WriteTimeout:   config.Listener.WriteTimeout,
		MaxHeaderBytes: config.Listener.MaxHeaderBytes,
	}

	proxyServer.Start()
}
