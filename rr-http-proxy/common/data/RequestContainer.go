package data

import (
	"net/http"
	"net/url"
)

type RequestContainer struct {
	RawURI    string
	Host      string
	ParsedURL url.URL
	Method    string
	Headers   http.Header
	Body      []byte
}
