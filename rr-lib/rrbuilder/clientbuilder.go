package rrbuilder

import (
	"github.com/CodeWizard-IL/rr-go/rr-lib/rrclient"
	"github.com/mitchellh/mapstructure"
)

type ClientConfig struct {
	Backend       BackendConfig
	Type          string
	Configuration map[string]any
}

type UnsupportedClientTypeError struct {
	ClientType string
}

func (e UnsupportedClientTypeError) Error() string {
	return "Unsupported client type: " + e.ClientType
}

func ClientFromConfig(config ClientConfig) (rrclient.RequestResponseClient, error) {
	switch config.Type {
	case "simple":
		backend, err := BackendFromConfig(config.Backend)
		if err != nil {
			return nil, err
		}
		client := rrclient.SimpleRequestResponseClient{Backend: backend}
		err = mapstructure.Decode(config.Configuration, &client)
		return &client, nil
	default:
		return nil, UnsupportedClientTypeError{config.Type}
	}
}
