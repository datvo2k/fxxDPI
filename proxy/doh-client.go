package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/sync/semaphore"
)

type dohClient struct {
	urlObject               url.URL
	sepaphoreAcquireTimeout time.Duration
	requestTimeout          time.Duration
	semaphore               *semaphore.Weighted
	dohJSONConverter        *dohJSONConverter
}
