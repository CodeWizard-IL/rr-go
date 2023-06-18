package main

import (
	"fmt"
	"forward/server"
	"gopkg.in/yaml.v3"
	"os"
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

	configFile, err := os.Open("config.yaml")
	if err != nil {
		fmt.Println("Error opening config file: ", err)
		return
	}
	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	var config ForwardProxyConfig
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file: ", err)
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
