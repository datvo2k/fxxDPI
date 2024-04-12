package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type dohJSONResponseAnswer struct {
	// Expires string `json:"Expires"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
	TTL     int    `json:"TTL"`
	Data    string `json:"data"`
}

type dohJSONResponse struct {
	Status int                     `json:"Status"`
	Answer []dohJSONResponseAnswer `json:"Answer"`
}
