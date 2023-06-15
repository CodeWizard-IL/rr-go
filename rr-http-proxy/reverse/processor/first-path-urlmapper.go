package processor

import (
	"fmt"
	"regexp"
)

type FirstPathUrlMapper struct {
	re *regexp.Regexp
}

func (mapper *FirstPathUrlMapper) MapURL(_ string, url string) string {
	if mapper.re == nil {
		mapper.re = regexp.MustCompile(`/([^/]*)/`)
	}

	matches := mapper.re.FindStringSubmatch(url)

	if len(matches) <= 1 {
		return ""
	}

	return fmt.Sprintf("http://%s%s", matches[1], url[len(matches[0])-1:])
}
