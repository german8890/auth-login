package http_client

import (
	"net/http"
)

type ClientOption func(*Client)

type Client interface {
	WithPreconfiguredClient(*http.Client) Client
	WithClientBefore(before ...RequestFunc) Client
	WithClientAfter(after ...ClientResponseFunc) Client
	WithClientFinalizer(f ...ClientFinalizerFunc) Client
	WithRetries(attempts int) Client
	Build() ClientS
}

type ClientS struct {
	C         *http.Client
	Attempts  int
	Before    []RequestFunc
	After     []ClientResponseFunc
	Finalizer []ClientFinalizerFunc
}

func New() Client {
	return new(ClientS)
}

func (c ClientS) WithPreconfiguredClient(h *http.Client) Client {
	c.C = h
	return c
}

func (c ClientS) WithClientBefore(before ...RequestFunc) Client {
	c.Before = append(c.Before, before...)
	return c
}

func (c ClientS) WithClientAfter(after ...ClientResponseFunc) Client {
	c.After = append(c.After, after...)
	return c
}

func (c ClientS) WithClientFinalizer(finalizer ...ClientFinalizerFunc) Client {
	c.Finalizer = append(c.Finalizer, finalizer...)
	return c
}

func (c ClientS) WithRetries(attempts int) Client {
	c.Attempts = attempts
	return c
}

func (c ClientS) Build() ClientS {

	if c.C == nil {
		c.C = http.DefaultClient
	}
	return c
}
