package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

type connections struct {
	sync.RWMutex
	lookup *ristretto.Cache
}

func NewConnection(lookup *ristretto.Cache) *connections {
	return &connections{lookup: lookup}
}

func (c *connections) Set(id uint16, addr *net.UDPAddr) {
	ttl := 10 * time.Second
	c.lookup.SetWithTTL(uint64(id), addr, 1, ttl)
	return
}

func (c *connections) Get(id uint16) *net.UDPAddr {
	val, ok := c.lookup.Get(uint64(id))
	if ok {
		return val.(*net.UDPAddr)
	}

	log.Printf("id not found, %v", id)

	return nil
}
