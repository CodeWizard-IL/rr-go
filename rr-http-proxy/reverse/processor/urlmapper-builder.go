package processor

import "github.com/mitchellh/mapstructure"

type UrlMapperConfig struct {
	Type          string
	Configuration map[string]any
}

type UnsupportedUrlMapperTypeError struct {
	UrlMapperType string
}

func (e UnsupportedUrlMapperTypeError) Error() string {
	return "Unsupported url mapper type: " + e.UrlMapperType
}

func UrlMapperFromConfig(config UrlMapperConfig) (URLMapper, error) {
	switch config.Type {
	case "as-is":
		urlMapper := AsIsUrlMapper{}
		err := mapstructure.Decode(config.Configuration, &urlMapper)
		if err != nil {
			return nil, err
		}
		return &urlMapper, nil

	case "first-path":
		urlMapper := FirstPathUrlMapper{}
		err := mapstructure.Decode(config.Configuration, &urlMapper)
		if err != nil {
			return nil, err
		}
		return &urlMapper, nil

	case "default-host":
		urlMapper := DefaultHostURLMapper{}
		err := mapstructure.Decode(config.Configuration, &urlMapper)
		if err != nil {
			return nil, err
		}
		return &urlMapper, nil
	default:
		return nil, UnsupportedUrlMapperTypeError{config.Type}
	}
}
