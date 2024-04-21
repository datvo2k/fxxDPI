package main

import (
	"net"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/elazarl/goproxy"
)

var httpClientTimeout = 15 * time.Second
var dialTimeout = 7 * time.Second

func NewHTTPReturn(configuration *Configuration) *fasthttp.Client {
	address := net.JoinHostPort(configuration.Host, configuration.Port)
	return &fasthttp.Client{
		ReadTimeout:         30 * time.Second,
		MaxConnsPerHost:     233,
		MaxIdleConnDuration: 15 * time.Minute,
		ReadBufferSize:      8192,
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialDualStackTimeout(address, dialTimeout)
		},
	}
}

func getDNSServer(configuration *Configuration) string {
	return configuration.DNSConfig.Server
}

func NewResolver(provider string) (*net.Resolver, error) {
	
}
