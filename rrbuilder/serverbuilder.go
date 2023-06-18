package rrbuilder

import (
	"github.com/mitchellh/mapstructure"
	"rrserver"
)

type ServerConfig struct {
	Backend       BackendConfig
	Type          string
	Configuration map[string]any
}

type UnsupportedServerTypeError struct {
	ServerType string
}

func (e UnsupportedServerTypeError) Error() string {
	return "Unsupported server type: " + e.ServerType
}

func ServerFromConfig(config ServerConfig) (rrserver.RequestResponseServer, error) {
	switch config.Type {
	case "simple":
		backend, err := BackendFromConfig(config.Backend)
		if err != nil {
			return nil, err
		}
		server := rrserver.SimpleRequestResponseServer{Backend: backend}
		err = mapstructure.Decode(config.Configuration, &server)
		return &server, nil
	default:
		return nil, UnsupportedServerTypeError{config.Type}
	}
}
