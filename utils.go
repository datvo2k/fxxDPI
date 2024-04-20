package main

import (
	"log"
	"time"
	"github.com/elazarl/goproxy"
)

var httpClientTimeout = 15 * time.Second
var dialTimeout = 7 * time.Second

