package processor

type URLMapper interface {
	MapURL(host string, url string) string
}
