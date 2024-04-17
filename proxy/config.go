package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type MetricsConfiguration struct {
	TimerIntervalSeconds int `json:"timerIntervalSeconds"`
}

// DNSServerConfiguration is the DNS server configuration.
type DNSServerConfiguration struct {
	ListenAddress HostAndPort `json:"listenAddress"`
}

// HostAndPort is a host and port.
type HostAndPort struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// DOHClientConfiguration is the DOH client configuration
type DOHClientConfiguration struct {
	URL                                 string `json:"url"`
	MaxConcurrentRequests               int64  `json:"maxConcurrentRequests"`
	SemaphoreAcquireTimeoutMilliseconds int    `json:"semaphoreAcquireTimeoutMilliseconds"`
	RequestTimeoutMilliseconds          int    `json:"requestTimeoutMilliseconds"`
}

// DNSProxyConfiguration is the proxy configuration.
type DNSProxyConfiguration struct {
	BlockedDomainsFile string `json:"blockedDomainsFile"`
	ClampMinTTLSeconds uint32 `json:"clampMinTTLSeconds"`
	ClampMaxTTLSeconds uint32 `json:"clampMaxTTLSeconds"`
}

// CacheConfiguration is the cache configuration.
type CacheConfiguration struct {
	MaxSize              int `json:"maxSize"`
	MaxPurgesPerTimerPop int `json:"maxPurgesPerTimerPop"`
	TimerIntervalSeconds int `json:"timerIntervalSeconds"`
}

// PrefetchConfiguration is the prefetch configuration.
type PrefetchConfiguration struct {
	MaxCacheSize            int `json:"maxCacheSize"`
	NumWorkers              int `json:"numWorkers"`
	SleepIntervalSeconds    int `json:"sleepIntervalSeconds"`
	MaxCacheEntryAgeSeconds int `json:"maxCacheEntryAgeSeconds"`
}

// PprofConfiguration is the pprof configuration.
type PprofConfiguration struct {
	ListenAddress string `json:"listenAddress"`
	Enabled       bool   `json:"enabled"`
}

// Configuration is the DNS proxy configuration.
type Configuration struct {
	MetricsConfiguration   MetricsConfiguration   `json:"metricsConfiguration"`
	DNSServerConfiguration DNSServerConfiguration `json:"dnsServerConfiguration"`
	DOHClientConfiguration DOHClientConfiguration `json:"dohClientConfiguration"`
	DNSProxyConfiguration  DNSProxyConfiguration  `json:"dnsProxyConfiguration"`
	PprofConfiguration     PprofConfiguration     `json:"pprofConfiguration"`
}

// ReadConfiguration reads the DNS proxy configuration from a json file.
func ReadConfiguration(configFile string) (*Configuration, error) {
	log.Printf("reading config file %q", configFile)

	source, err := os.ReadFile(configFile)
	if err != nil {
		err = fmt.Errorf("ioutil.ReadFile error: %w", err)
		return nil, err
	}

	var config Configuration
	if err = json.Unmarshal(source, &config); err != nil {
		err = fmt.Errorf("json.Unmarshal error: %w", err)
		return nil, err
	}

	return &config, nil
}

// JoinHostPort joins the host and port.
func (hostAndPort *HostAndPort) joinHostPort() string {
	return net.JoinHostPort(hostAndPort.Host, hostAndPort.Port)
}
