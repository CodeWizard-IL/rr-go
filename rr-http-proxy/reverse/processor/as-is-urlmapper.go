package processor

import "fmt"

type AsIsUrlMapper struct {
}

func (mapper *AsIsUrlMapper) MapUrl(host string, url string) string {
	return fmt.Sprintf("http://%s%s", host, url)
}
