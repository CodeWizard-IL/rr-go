package processor

type DefaultHostURLMapper struct {
	DefaultHost string
}

func (mapper *DefaultHostURLMapper) MapURL(_ string, url string) string {
	return "http://" + mapper.DefaultHost + url
}
