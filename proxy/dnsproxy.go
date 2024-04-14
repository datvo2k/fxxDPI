package proxy

import (
	"context"
	"github.com/miekg/dns"
	"log"
)

// DNSProxy is the DNS proxy.
type DNSProxy interface {
	Start()
}

type dnsProxy struct {
	configuration *Configuration
	metrics       *metrics
	dnsServer     *dnsServer
	dohClient     *dohClient
}

func NewDNSProxy(configuration *Configuration) DNSProxy {
	metrics := newMetrics(&configuration.MetricsConfiguration)

	return &dnsProxy{
		configuration: configuration,
		metrics:       metrics,
		dnsServer:     newDNSServer(&configuration.DNSServerConfiguration),
		dohClient:     newDOHClient(configuration.DOHClientConfiguration, newDOHJSONConverter(metrics)),
	}
}

func (dnsProxy *dnsProxy) clampAndGetMinTTLSeconds(m *dns.Msg) uint32 {
	clampMinTTLSeconds := dnsProxy.configuration.DNSProxyConfiguration.ClampMinTTLSeconds
	clampMaxTTLSeconds := dnsProxy.configuration.DNSProxyConfiguration.ClampMaxTTLSeconds

	foundRRHeaderTTL := false
	rrHeaderMinTTLSeconds := clampMinTTLSeconds

	processRRHeader := func(rrHeader *dns.RR_Header) {
		ttl := rrHeader.Ttl
		if ttl < clampMinTTLSeconds {
			ttl = clampMinTTLSeconds
		}
		if ttl > clampMaxTTLSeconds {
			ttl = clampMaxTTLSeconds
		}
		if (!foundRRHeaderTTL) || (ttl < rrHeaderMinTTLSeconds) {
			rrHeaderMinTTLSeconds = ttl
			foundRRHeaderTTL = true
		}
		rrHeader.Ttl = ttl
	}

	for _, rr := range m.Answer {
		processRRHeader(rr.Header())
	}
	for _, rr := range m.Ns {
		processRRHeader(rr.Header())
	}
	for _, rr := range m.Extra {
		rrHeader := rr.Header()
		if rrHeader.Rrtype != dns.TypeOPT {
			processRRHeader(rrHeader)
		}
	}

	return rrHeaderMinTTLSeconds
}

func (dnsProxy *dnsProxy) writeResponse(w dns.ResponseWriter, response *dns.Msg) {
	if err := w.WriteMsg(response); err != nil {
		dnsProxy.metrics.incrementWriteResponseErrors()
		log.Printf("writeResponse error = %v", err)
	}
}

func (dnsProxy *dnsProxy) createProxyHandlerFunc() dns.HandlerFunc {
	return func(w dns.ResponseWriter, request *dns.Msg) {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if len(request.Question) != 1 {
			log.Printf("bad request.Question length %v request %v", len(request.Question), request)
			dns.HandleFailed(w, request)
			return
		}

		requestID := request.Id

		request.Id = 0
		responseMsg, err := dnsProxy.dohClient.makeRequest(ctx, request)
		if err != nil {
			dnsProxy.metrics.incrementDOHClientErrors()
			log.Printf("makeHttpRequest error: %v", err)
			request.Id = requestID
			dns.HandleFailed(w, request)
			return
		}

		responseMsg.Id = requestID
		dnsProxy.writeResponse(w, responseMsg)
	}
}

func (dnsProxy *dnsProxy) createBlockedDomainHandlerFunc() dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		dnsProxy.metrics.incrementBlocked()

		responseMsg := new(dns.Msg)
		responseMsg.SetRcode(r, dns.RcodeNameError)
		dnsProxy.writeResponse(w, responseMsg)
	}
}

func (dnsProxy *dnsProxy) createServeMux() *dns.ServeMux {

	dnsServeMux := dns.NewServeMux()
	dnsServeMux.HandleFunc(".", dnsProxy.createProxyHandlerFunc())

	if len(dnsProxy.configuration.DNSProxyConfiguration.BlockedDomainsFile) > 0 {
		blockedHandler := dnsProxy.createBlockedDomainHandlerFunc()
		installHandlersForBlockedDomains(dnsProxy.configuration.DNSProxyConfiguration.BlockedDomainsFile, dnsServeMux, blockedHandler)
	}

	return dnsServeMux
}

func (dnsProxy *dnsProxy) Start() {
	log.Printf("begin dnsProxy.Start")
	dnsProxy.metrics.start()
	dnsProxy.dnsServer.start(dnsProxy.createServeMux())
	startPprof(&dnsProxy.configuration.PprofConfiguration)
	log.Printf("end dnsProxy.Start")
}
