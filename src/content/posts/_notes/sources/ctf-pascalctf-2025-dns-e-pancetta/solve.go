package main

import (
	"encoding/hex"
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	handle, err := pcap.OpenOffline("./capture.pcapng")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	var sequence string

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		dnsLayer := packet.Layer(layers.LayerTypeDNS)
		if dnsLayer == nil {
			continue
		}

		dns, _ := dnsLayer.(*layers.DNS)
		if !dns.QR {
			continue
		}

		host := string(dns.Questions[0].Name)
		log.Printf("host: %s", host)

		if !strings.HasSuffix(host, ".attacker.com") {
			continue
		}

		host = strings.TrimSuffix(host, ".attacker.com")
		sequence += host
	}

	answer, _ := hex.DecodeString(sequence)

	log.Printf("answer: %s", answer)
}
