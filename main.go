package main

import (
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	log.Print("Starting")
	origin := &net.UDPAddr{Port: 53}
	l, err := net.ListenUDP("udp", origin)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	defer l.Close()
	// message := map[uint16]*net.UDPAddr{}
	message := NewConnection()

	for {
		buf := make([]byte, 512)
		_, addr, err := l.ReadFromUDP(buf)
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go func(message *connections) {
			var m dnsmessage.Message
			err = m.Unpack(buf)
			if err != nil {
				log.Fatal("unpack err:", err)
			}

			question := m.Questions[0]
			if strings.Contains(question.Name.String(), "ec2.internal") {
				message.Set(m.Header.ID, addr)

				re := regexp.MustCompile("(\\d{1,3})-(\\d{1,3})-(\\d{1,3})-(\\d{1,3})")
				toCheck := []byte(question.Name.String())
				matches := re.FindSubmatch(toCheck)

				resolved := [4]byte{}
				for i, v := range matches[1:] {
					str := string(v)
					val, _ := strconv.Atoi(str)
					resolved[i] = uint8(val)
				}

				m.Header.Response = true
				m.Header.Authoritative = true
				newAnswers := []dnsmessage.Resource{
					{
						Header: dnsmessage.ResourceHeader{
							Name:   dnsmessage.MustNewName(m.Questions[0].Name.String()),
							Type:   dnsmessage.TypeA,
							Class:  dnsmessage.ClassINET,
							TTL:    278,
							Length: 22,
						},

						Body: &dnsmessage.AResource{A: resolved},
					},
				}

				m.Answers = newAnswers
			}

			packed, err := m.Pack()
			if err != nil {
				log.Printf("Packing err: %s", err)
			}
			if m.Header.Response {
				_, err := l.WriteToUDP(packed, message.Get(m.Header.ID))
				if err != nil {
					log.Printf("Write err: %s", err)
				}
				if len(m.Answers) > 0 {
					log.Printf("Questions: %v, type: %v, Answers: %v", m.Questions, m.Answers[0].Header.GoString(), m.Answers[0].Body.GoString())
				}
			} else {
				message.Set(m.Header.ID, addr)
				resolver := net.UDPAddr{IP: net.IP{1, 1, 1, 1}, Port: 53}
				_, err = l.WriteToUDP(packed, &resolver)
				if err != nil {
					log.Printf("failed to resolve %s", err)
				}
			}
		}(message)
	}

	log.Println("exiting")
}
