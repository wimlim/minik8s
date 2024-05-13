package dnsserver

import (
	"fmt"
	"log"
	"minik8s/pkg/etcd"
	"strings"

	"github.com/miekg/dns"
)

func addDNSRecord(domain, ip string) {
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}
	if err := etcd.EtcdKV.Put(domain, []byte(ip)); err != nil {
		log.Fatalf("Failed to add DNS record: %v", err)
	} else {
		log.Printf("DNS record added: %s -> %s", domain, ip)
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		fmt.Println("Query")
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		ip, err := etcd.EtcdKV.Get(q.Name)
		if err == nil && ip != nil {
			rrString := fmt.Sprintf("%s 3600 IN A %s", q.Name, strings.TrimSpace(string(ip)))
			rr, err := dns.NewRR(rrString)
			if err == nil {
				m.Answer = append(m.Answer, rr)
			} else {
				log.Printf("Failed to create DNS RR: %v", err)
			}
		} else {
			log.Printf("Failed to retrieve IP for %s: %v", q.Name, err)
		}
	}
}

func Run() {
	addDNSRecord("minik8s.com", "1.2.2.3")

	server := &dns.Server{Addr: ":10053", Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)
	log.Println("Starting DNS server...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
