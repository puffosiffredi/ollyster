package irc

import (
	"log"
	"net"
	"strings"
)

func (this *IrcServer) ReadIpFromHost() *net.TCPAddr {

	var addr string

	// first we get our IP
	conn, err := net.Dial("udp", this.serveraddr+":"+this.serverport)
	if err != nil {
		log.Printf("[DNS] SYSADMIIIIIN : cannot use UDP")
		conn.Close()
		return nil

	} else {
		addr, _, _ = net.SplitHostPort(conn.LocalAddr().String())
		log.Printf("[DNS] Local addr string: %s", addr)
		conn.Close()
	}

	// now get a free TCP Port

	list, _ := net.Listen("tcp", ":0")
	_, port, _ := net.SplitHostPort(list.Addr().String())
	list.Close()

	// then put all together

	ind, _ := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, port))

	log.Printf("[DNS] Resolution ok: %s -> %s", "Local", ind.IP)

	return ind
}

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
