package proxy

import (
	"net"
	"sync"
)

var cacheAddrMapLock sync.RWMutex
var cacheTCPAddrMap = map[string]*net.TCPAddr{}
var domainProxiesCache = map[string]bool{}
var domainProxiesCacheLock sync.RWMutex
