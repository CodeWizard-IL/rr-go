package rrbuilder

import (
	"github.com/mitchellh/mapstructure"
	"rrbackend"
	"rrbackend/local"
	"rrbackendamqp09"
	"rrbackendazsmb"
)

type UnsupportedBackendTypeError struct {
	BackendType string
}

func (e UnsupportedBackendTypeError) Error() string {
	return "Unsupported backend type: " + e.BackendType
}

type BackendConfig struct {
	Type          string
	Configuration map[string]any
}

func BackendFromConfig(config BackendConfig) (rrbackend.RequestResponseBackend, error) {
	switch config.Type {
	case "amqp09":
		amqp09Backend := rrbackendamqp09.Amqp09Backend{}
		err := mapstructure.Decode(config.Configuration, &amqp09Backend)
		if err != nil {
			return nil, err
		}
		return &amqp09Backend, nil
	case "local":
		return &local.RequestResponseBackend{}, nil
	case "azsmb":
		rrBackendAzSMB := rrbackendazsmb.RRBackendAzSMB{}
		err := mapstructure.Decode(config.Configuration, &rrBackendAzSMB)
		if err != nil {
			return nil, err
		}
		return &rrBackendAzSMB, nil

	default:
		return nil, UnsupportedBackendTypeError{config.Type}
	}
}
