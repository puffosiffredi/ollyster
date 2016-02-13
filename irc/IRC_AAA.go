package irc

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

func init() {

	go irc_AAA()

}

func irc_AAA() {

	conn, err := net.Dial("tcp4", "130.239.18.119:8000")
	if err != nil {
		log.Printf("[AAA] CONNECTION ERROR: %s", err)
	} else {
		log.Println("[AAA] Connected. Now waiting 5 seconds")
		time.Sleep(3000 * time.Millisecond)
		log.Println("[AAA] Now sending the AAA")
		conn.Write([]byte("CAP LS\n"))
		time.Sleep(3000 * time.Millisecond)
		conn.Write([]byte("NICK Ollyster\n"))
		time.Sleep(3000 * time.Millisecond)
		conn.Write([]byte("USER ollyster ollyster irc.oftc.net :Ollyster\n"))
		conn.Write([]byte("CAP END\n"))
		time.Sleep(3000 * time.Millisecond)
		conn.Write([]byte("JOIN #social \n"))
		time.Sleep(3000 * time.Millisecond)
		log.Println("[AAA] AAA terminated, now listening")
		linea := make([]byte, 1024)

		for {

			linea, _, _ = bufio.NewReader(conn).ReadLine()

			message := string(linea)

			if matches, _ := regexp.MatchString("(?i)^PING :.*$", message); matches == true {
				log.Printf("[IRC] %s ", message)
				sinta := strings.Split(message, ":")
				conn.Write([]byte("PONG :" + sinta[1]))
				log.Printf("[IRC] Sending back the -> %s", "PONG :"+sinta[1])
				continue
			}

			// :nick!user@ip-address PRIVMSG your-nick :VERSION
			if matches, _ := regexp.MatchString("(?i)^:.*PRIVMSG.*VERSION$", message); matches == true {
				sinta := strings.Split(message, "!")
				log.Printf("[IRC] VERSION REQUEST from %s ", sinta[0])
				version := "NOTICE " + strings.TrimLeft(sinta[0], ":") + " : VERSION Ollyster DEV https://github.com/uriel-fanelli/ollyster"
				conn.Write([]byte(version))
				log.Printf("[IRC] Sending back the -> %s", version)
				continue
			}

			// This must be the last one
			if len(message) > 1 {
				log.Printf("[IRC]>  " + message)
				continue
			}

		}

	}

}

func IrcInitialize() {
	log.Println("[AAA] Initializing IRC Engine...")
}
