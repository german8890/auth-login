package client

import (
	"net/http"
	"time"
)

type Client interface {
	Timeout(timeOut time.Duration)
	IdleConnTimeout(timeOut time.Duration)
	TLSHandshakeTimeout(timeOut time.Duration)
	ExpectContinueTimeout(timeOut time.Duration)
	MaxIdleConns(maxIdleConns int)
	MaxConnsPerHost(maxConnsPerHost int)
	MaxIdleConnsPerHost(maxIdleConnsPerHost int)
	DisableKeepAlives(boolean bool)
	Build() *http.Client
}

type httpClient struct {
	timeOut               time.Duration
	maxIdleConns          int
	maxConnsPerHost       int
	maxIdleConnsPerHost   int
	idleConnTimeout       time.Duration
	tLSHandshakeTimeout   time.Duration
	expectContinueTimeout time.Duration
	disableKeepAlives     bool
}

func MakeHttpClient() Client {
	return &httpClient{}
}

func (c *httpClient) Build() *http.Client {
	client := &http.Client{}
	transport := http.DefaultTransport.(*http.Transport).Clone()

	if c.timeOut != 0 {
		client.Timeout = c.timeOut * time.Millisecond
	}

	if c.idleConnTimeout != 0 {
		transport.IdleConnTimeout = c.idleConnTimeout * time.Millisecond
	}
	if c.tLSHandshakeTimeout != 0 {
		transport.TLSHandshakeTimeout = c.tLSHandshakeTimeout * time.Millisecond
	}
	if c.expectContinueTimeout != 0 {
		transport.ExpectContinueTimeout = c.expectContinueTimeout * time.Millisecond
	}
	if c.maxConnsPerHost != 0 {
		transport.MaxConnsPerHost = c.maxConnsPerHost
	}
	if c.maxIdleConnsPerHost != 0 {
		transport.MaxIdleConnsPerHost = c.maxIdleConnsPerHost
	}
	if c.disableKeepAlives {
		transport.DisableKeepAlives = true
	}

	client.Transport = transport

	return client

}

func (c *httpClient) Timeout(timeOut time.Duration) {
	c.timeOut = timeOut
}

func (c *httpClient) DisableKeepAlives(boolean bool) {
	c.disableKeepAlives = boolean
}

func (c *httpClient) IdleConnTimeout(timeOut time.Duration) {
	c.idleConnTimeout = timeOut
}

func (c *httpClient) TLSHandshakeTimeout(timeOut time.Duration) {
	c.tLSHandshakeTimeout = timeOut
}

func (c *httpClient) ExpectContinueTimeout(timeOut time.Duration) {
	c.expectContinueTimeout = timeOut
}

func (c *httpClient) MaxIdleConns(maxIdleConns int) {
	c.maxIdleConns = maxIdleConns
}

func (c *httpClient) MaxConnsPerHost(maxConnsPerHost int) {
	c.maxConnsPerHost = maxConnsPerHost
}

func (c *httpClient) MaxIdleConnsPerHost(maxIdleConnsPerHost int) {
	c.maxIdleConnsPerHost = maxIdleConnsPerHost
}
