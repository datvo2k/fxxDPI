package proxy

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/AdguardTeam/dnsproxy/proxy"
	"github.com/miekg/dns"
)

// Config is the DNS proxy configuration.
type Config struct {
	// ListenAddr is the address the DNS server is supposed to listen to.
	ListenAddr netip.AddrPort

	// Upstream is the upstream that the requests will be forwarded to.  The
	// format of an upstream is the one that can be consumed by
	// [proxy.ParseUpstreamsConfig].
	Upstream string

	// RedirectIPv4To is the IP address A queries will be redirected to.
	RedirectIPv4To net.IP

	// RedirectIPv6To is the IP address AAAA queries will be redirected to.
	RedirectIPv6To net.IP
}

type DNSProxy struct {
	proxy          *proxy.Proxy
	redirectIPv4To net.IP
	redirectIPv6To net.IP
}

// defaultTTL is the default TTL for the rewritten records.
const defaultTTL = 60

// type check
// var _ io.Closer = (*DNSProxy)(nil)

// New creates a new instance of *DNSProxy
func New(cfg *Config) (d *DNSProxy, err error) {
	proxyConfig, err := createProxyConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("dnsproxy: invalid configuration: %w", err)
	}

	d = &DNSProxy{
		redirectIPv4To: cfg.RedirectIPv4To,
		redirectIPv6To: cfg.RedirectIPv6To,
	}
	d.proxy = &proxy.Proxy{
		Config: proxyConfig,
	}
	d.proxy.RequestHandler = d.requestHandler

	return d, nil

}

// createProxyConfig creates DNS proxy configuration.
func createProxyConfig(cfg *Config) (proxyConfig proxy.Config, err error) {
	upstreamCfg, err := proxy.ParseUpstreamsConfig([]string{cfg.Upstream}, nil)
	if err != nil {
		return proxyConfig, fmt.Errorf("failed to parse upstream %s: %w", cfg.Upstream, err)
	}
	ip := net.IP(cfg.ListenAddr.Addr().AsSlice())

	udpPort := &net.UDPAddr{
		IP:   ip,
		Port: int(cfg.ListenAddr.Port()),
	}
	tcpPort := &net.TCPAddr{
		IP:   ip,
		Port: int(cfg.ListenAddr.Port()),
	}

	proxyConfig.UDPListenAddr = []*net.UDPAddr{udpPort}
	proxyConfig.TCPListenAddr = []*net.TCPAddr{tcpPort}
	proxyConfig.UpstreamConfig = upstreamCfg

	return proxyConfig, nil
}

func (d *DNSProxy) requestHandler(p *proxy.Proxy, ctx *proxy.DNSContext) (err error) {
	qType := ctx.Req.Question[0].Qtype
	if qType != dns.TypeA && qType != dns.TypeAAAA {
		// Doing nothing with the request if it's not A/AAAA, we cannot
		// rewrite them anyway.
		return nil
	}

	return p.Resolve(ctx)
}

// Start starts the DNSProxy server.
// func (d *DNSProxy) Start() (err error) {
// 	err = d.proxy.Start()
// 	return err
// }

// Close implements the [io.Closer] interface for DNSProxy.
// func (d *DNSProxy) Close() (err error) {
// 	err = d.proxy.Stop()
// 	return err
// }
