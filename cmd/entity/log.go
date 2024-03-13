package entity

import "net/http"

type Logger struct {
	Headers       http.Header `json:"headers"`
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	UriPath       string      `json:"uri_path"`
	Body          interface{} `json:"body"`
	BodyFormatted interface{} `json:"bodyFormatted,omitempty"`
	BodyError     interface{} `json:"bodyError,omitempty"`
}

type Service struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}
