package irc

import (
	"bufio"
	"log"
	"net"
	"ollyster/conf"
	"ollyster/files"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//this is to contain all structs we need.
//will be populated using ./etc
type IrcServer struct {
	servername string
	serverport string
	serveraddr string
	nickname   string
	socket     *net.TCPConn
	protocol   string // can be tcp4, tcp6 , tcp (automatic detect)
	delay      int    // milliseconds to wait when sending Dial commands
	heartbeat  int    // keepalive in seconds
	channel    string // the channel to join
}

func init() {

	var MyServer IrcServer

	// taken by the conf file
	MyServer.servername = conf.OConfig["servername"]
	MyServer.serverport = conf.OConfig["serverport"]
	MyServer.serveraddr = conf.OConfig["serveraddr"]
	MyServer.nickname = conf.OConfig["nickname"]
	MyServer.protocol = conf.OConfig["protocol"]
	MyServer.delay, _ = strconv.Atoi(conf.OConfig["delay"])
	MyServer.heartbeat, _ = strconv.Atoi(conf.OConfig["heartbeat"])
	MyServer.channel = conf.OConfig["channel"]

	go MyServer.ircClient()

}

func (this *IrcServer) ircClient() {

	this.ircDial()

	linea := make([]byte, 1024)

	for {

		var err error
		var exceed bool

		reader := bufio.NewReader(this.socket)

		linea, exceed, err = reader.ReadLine()

		if err != nil {
			this.socket.Close()
			this.ircDial()
			continue
		}

		message := string(linea)

		if matches, _ := regexp.MatchString("(?i)^PING :.*$", message); matches == true {
			log.Printf("[IRC] %s ", message)
			sinta := strings.Split(message, ":")
			this.socket.Write([]byte("PONG :" + sinta[1]))
			log.Printf("[IRC] Sending back the -> %s", "PONG :"+sinta[1])
			continue
		}

		// :nick!user@ip-address PRIVMSG your-nick :VERSION
		if matches, _ := regexp.MatchString("(?i)^:.*PRIVMSG.*VERSION$", message); matches == true {
			sinta := strings.Split(message, "!")
			log.Printf("[IRC] VERSION REQUEST from %s ", sinta[0])
			version := "NOTICE " + strings.TrimLeft(sinta[0], ":") + " : VERSION Ollyster DEV https://github.com/uriel-fanelli/ollyster"
			this.socket.Write([]byte(version))
			log.Printf("[IRC] Sending back the -> %s", version)
			continue
		}

		// :nick!user@ip-address PRIVMSG your-nick :Message
		privMsgString := "(?i)^:.*!.*PRIVMSG.*" + this.nickname + " :.*$"
		if matches, _ := regexp.MatchString(privMsgString, message); matches == true {
			payload := strings.TrimLeft(message, ":") // starting a payload with a separator? WTF?
			sinta := strings.Split(message, "!")
			sender := strings.TrimLeft(sinta[0], ":")
			// sender contains the sender nick

			torn := strings.Split(payload, ":")
			msg := strings.Join(torn[1:], ":")
			// msg contains the message after the 1st colon

			files.MyStream.WriteMsgPriv(sender, msg)
			log.Printf("[IRC] Private message from %s:  <%s>", sender, msg)
			continue
		}

		// :nick!user@ip-address PRIVMSG #channel :Message
		chanMsgString := "(?i)^:.*!.*PRIVMSG.*" + this.channel + " :.*$"
		if matches, _ := regexp.MatchString(chanMsgString, message); matches == true {
			payload := strings.TrimLeft(message, ":") // starting a payload with a separator? WTF?
			sinta := strings.Split(message, "!")
			sender := strings.TrimLeft(sinta[0], ":")
			// sender contains the sender nick

			torn := strings.Split(payload, ":")
			msg := strings.Join(torn[1:], ":")
			// msg contains the message after the 1st colon

			log.Printf("[IRC] %s sent a message to %s:  <%s>", sender, this.channel, msg)
			files.MyStream.WriteMsgGroup(sender, this.channel, msg)
			continue
		}

		// This must be the last one
		if err == nil {
			log.Printf("[IRC](exceed = %t)>  %s ", exceed, message)
			continue
		}

	}

}

func (this *IrcServer) ircDial() {

	var err error

	this.socket, err = net.DialTCP("tcp4", nil, this.MakeAddr())

	if err != nil {
		log.Printf("[AAA] CONNECTION ERROR: %s", err)
	} else {
		this.socket.SetKeepAlive(true)                                              // keepalive on
		this.socket.SetKeepAlivePeriod(time.Duration(this.heartbeat) * time.Second) // tcp keepalive to 10 seconds
		this.socket.SetLinger(0)                                                    // brutal close

		log.Println("[AAA] Connected. Now waiting for courtesy")
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		log.Println("[AAA] Now sending the AAA")
		this.socket.Write([]byte("CAP LS\n"))
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		this.socket.Write([]byte("NICK " + this.nickname + "\n"))
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		userString := "USER " + this.nickname + " " + this.nickname + " " + this.servername + " :" + this.nickname + "\n"
		log.Print("[AAA] " + userString)
		this.socket.Write([]byte(userString))
		this.socket.Write([]byte("CAP END\n"))
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		log.Println("[AAA] AAA terminated, now joining")
		this.IrcJoin()

	}

}

func (this *IrcServer) MakeAddr() *net.TCPAddr {

	tcpAddr, err := net.ResolveTCPAddr("tcp", this.serveraddr+":"+this.serverport)
	if err != nil {

		log.Printf("[TCP] Error in creating TCP Address %s", err)
		return nil
	}

	return tcpAddr

}

func (this *IrcServer) IrcJoin() {

	joinString := "JOIN " + this.channel + "\n"

	log.Printf("[IRC] Joining  %s", this.channel)
	_, err := this.socket.Write([]byte(joinString))
	if err != nil {
		log.Printf("[IRC] Cannot join  %s: %s", this.channel, err)

	} else {
		log.Printf("[IRC] Successfully joined %s", this.channel)
	}

}

func IrcInitialize() {
	log.Println("[AAA] Initializing IRC Engine...")
}
