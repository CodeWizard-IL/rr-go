package util

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadYamlToStruct[T any](yamlFile string, config *T) error {
	configFile, err := os.Open(yamlFile)
	if err != nil {
		fmt.Println("Error opening config file: ", err)
		return err
	}
	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(config)
	if err != nil {
		fmt.Println("Error decoding config file: ", err)
		return err
	}
	return err
}
