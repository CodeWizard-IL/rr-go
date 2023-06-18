package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"reverse/processor"
	"reverse/server"
	"rrbuilder"
	"rrserver"
)

type ReverseProxyConfig struct {
	Server    rrbuilder.ServerConfig
	UrlMapper processor.UrlMapperConfig
}

func main() {
	fmt.Println("Request Response HTTP Proxy - Receiver")

	// Read ReverseProxyConfig from YAML file

	configFile, err := os.Open("config.yaml")
	if err != nil {
		fmt.Println("Error opening config file: ", err)
		return
	}
	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	var config ReverseProxyConfig
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config file: ", err)
		return
	}

	serverConfig := config.Server

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

	// Build the URLMapper

	urlMapper, err := processor.UrlMapperFromConfig(config.UrlMapper)
	if err != nil {
		fmt.Println("Error creating URLMapper: ", err)
		return
	}
	reverseProxyServer := server.ReverseProxyServer{
		RRServer:  *simpleServer,
		UrlMapper: urlMapper,
	}

	reverseProxyServer.Start()

	// Wait forever
	select {}

}
