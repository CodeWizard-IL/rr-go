package processor

import "fmt"

type AsIsUrlMapper struct {
}

func (mapper *AsIsUrlMapper) MapURL(host string, url string) string {
	return fmt.Sprintf("http://%s%s", host, url)
}
