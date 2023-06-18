package main

import (
	"common/util"
	"fmt"
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

	var config ReverseProxyConfig
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
