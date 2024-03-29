package main

import (
	"fmt"
	"github.com/CodeWizard-IL/rr-go/rr-lib/rrbuilder"
	"os"
	"rr-http-proxy/common/util"
	"rr-http-proxy/forward/server"
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

	const defaultConfigYaml = "config.yaml"

	// Get config YAML from environment variable
	configYaml := os.Getenv("CONFIG_YAML")
	if configYaml == "" {
		configYaml = defaultConfigYaml
	}

	err := util.ReadYamlToStruct(configYaml, &config)
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
