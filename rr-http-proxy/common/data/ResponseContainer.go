package data

import "net/http"

type ResponseContainer struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}
