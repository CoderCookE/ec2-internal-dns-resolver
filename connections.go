package main

import (
	"net"
	"sync"
)

type connections struct {
	sync.RWMutex
	Lookup map[uint16]*net.UDPAddr
}

func NewConnection() *connections {
	return &connections{Lookup: map[uint16]*net.UDPAddr{}}
}

func (c *connections) Set(id uint16, addr *net.UDPAddr) {
	c.Lock()
	defer c.Unlock()

	c.Lookup[id] = addr

	return
}

func (c *connections) Delete(id uint16) {
	c.Lock()
	defer c.Unlock()
	delete(c.Lookup, id)

	return
}

func (c *connections) Get(id uint16) *net.UDPAddr {
	c.Lock()
	defer c.Unlock()

	return c.Lookup[id]
}
