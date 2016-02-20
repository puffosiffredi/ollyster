package irc

import (
	"log"
	"net"
	"strings"
)

func (this *IrcServer) MakeAddr() *net.TCPAddr {

	tcpAddr, err := net.ResolveTCPAddr(this.protocol, this.serveraddr+":"+this.serverport)
	if err != nil {

		log.Printf("[TCP] Error in creating TCP Address %s", err)
		return nil
	}

	return tcpAddr

}

func (this *IrcServer) Resolve(fqdn string) (addr string) {

	addresses, err := net.LookupIP(fqdn)

	if err != nil {
		log.Println("[DNS] ERROR %s", err)
		this.protocol = "tcp4"
		return "127.0.0.1"

	} else {

		addr := addresses[0].String()

		log.Printf("[DNS] Resolution ok: %s -> %s", fqdn, addr)

		if strings.Contains(addr, ":") {
			addr = "[" + addr + "]"
			this.protocol = "tcp6"
			return addr
		}

		if strings.Contains(addr, ".") {
			this.protocol = "tcp4"
			return addr
		}

	}

	return "127.0.0.1"

}
