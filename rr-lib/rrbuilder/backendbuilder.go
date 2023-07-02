package rrbuilder

import (
	"github.com/mitchellh/mapstructure"
	"rr-lib/rrbackend"
	"rr-lib/rrbackend/amqp09"
	"rr-lib/rrbackend/azsb"
	"rr-lib/rrbackend/local"
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
		amqp09Backend := amqp09.RRBackendAmqp09{}
		err := mapstructure.Decode(config.Configuration, &amqp09Backend)
		if err != nil {
			return nil, err
		}
		return &amqp09Backend, nil
	case "local":
		return &local.RRBackendLocal{}, nil
	case "azsb":
		rrBackendAzSMB := azsb.RRBackendAzSB{}
		err := mapstructure.Decode(config.Configuration, &rrBackendAzSMB)
		if err != nil {
			return nil, err
		}
		return &rrBackendAzSMB, nil

	default:
		return nil, UnsupportedBackendTypeError{config.Type}
	}
}
