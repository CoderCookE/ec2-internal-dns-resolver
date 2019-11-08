package main

import (
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"net"
)

func main() {
	log.Print("Starting")
	origin := &net.UDPAddr{Port: 1053}
	l, err := net.ListenUDP("udp", origin)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()
	message := map[uint16]*net.UDPAddr{}

	for {
		buf := make([]byte, 512)
		_, addr, err := l.ReadFromUDP(buf)

		if err != nil {
			log.Fatal("accept error:", err)
		}

		var m dnsmessage.Message
		err = m.Unpack(buf)
		if err != nil {
			log.Fatal("unpack err:", err)
		}

		log.Printf("%s", m)
		log.Printf("%s", m.Header)
		log.Printf("%s", m.Header.Response)

		packed, err := m.Pack()

		if m.Header.Response {
			println("MADE IT HERE")
			log.Printf("%s", message[m.Header.ID])
			i, err := l.WriteToUDP(packed, message[m.Header.ID])
			delete(message, m.Header.ID)
			println(i)
			if err != nil {
				log.Fatal("Write err:", err)
			}
			if len(m.Answers) > 0 {
				log.Printf("Questions: %v, type: %v, Answers: %v", m.Questions, m.Answers[0].Header.GoString(), m.Answers[0].Body.GoString())
			}
		} else {
			message[m.Header.ID] = addr

			resolver := net.UDPAddr{IP: net.IP{1, 1, 1, 1}, Port: 53}
			_, err = l.WriteToUDP(packed, &resolver)
			if err != nil {
				log.Fatal("failed to resolve", err)
			}
		}
	}
	println("exiting")
}
