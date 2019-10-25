package main

import (
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"net"
)

func main() {
	log.Print("Starting")

	l, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	for {
		buf := make([]byte, 512)
		_, addr, err := l.ReadFromUDP(buf)
		go func() {
			if err != nil {
				log.Fatal("accept error:", err)
			}
			var m dnsmessage.Message
			err = m.Unpack(buf)
			if err != nil {
				log.Fatal("unpack err:", err)
			}

			// question := m.Questions[0]
			// log.Printf("%v", question.Name.String())
			// if strings.Contains(question.Name.String(), "facebook") {
			// 	println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&& facebook")
			// }

			packed, err := m.Pack()
			if err != nil {
				log.Fatal("Pack err:", err)
			}

			log.Printf("here")
			resolver := net.UDPAddr{IP: net.IP{1, 1, 1, 1}, Port: 53}
			_, err = l.WriteToUDP(packed, &resolver)

			if m.Header.Response {
				log.Printf("Questions: %v, Answers:", m.Questions, m.Answers[0].Body.GoString())
				packed, err = m.Pack()
				_, err = l.WriteToUDP(packed, addr)
				if err != nil {
					log.Fatal("Write err:", err)
				}
			}
		}()
	}
}
