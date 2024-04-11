package proxy

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// DNSProxy is the DNS proxy.
type DNSProxy interface {
	Start()
}

