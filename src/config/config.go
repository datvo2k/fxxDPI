package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Configuration represents a configuration object with properties for host, port, DNS configuration, and blocked domain file.
type Configuration struct {
	Host              string    `json:"host"`
	Port              string    `json:"port"`
	DNSConfig         DNSConfig `json:"dns"`
	BlockedDomainFile string    `json:"blockedDomainFile"`
}


// DNSConfig represents DNS configuration with properties for type, server, and cache size.
type DNSConfig struct {
	Type      string `json:"type"`
	Server    string `json:"server"`
	CacheSize uint32 `json:"cacheSize"`
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
