package main

import (
	"fmt"

	c "github.com/hashicorp/consul/api"
)

//GetConsulClient connects to consul and returns a consul client
func GetConsulClient(addr string) (*c.KV, error) {
	config := c.DefaultConfig()
	if len(addr) > 0 {
		config.Address = addr
	}
	consulClient, err := c.NewClient(config)
	if err != nil {
		return nil, err
	}
	return consulClient.KV(), nil
}

//ConsulGetKey get a raw key from consul and returns it
func ConsulGetKey(consul *c.KV, key string) ([]byte, error) {
	kvpair, _, err := consul.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kvpair == nil {
		return nil, fmt.Errorf("consul missing key: %v", key)
	}

	return kvpair.Value, nil
}
