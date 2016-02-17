package irc

import (
	"bufio"
	"log"
	"net"
	"ollyster/conf"
	"strconv"
	"time"
)

//this is to contain all structs we need.
//will be populated using ./etc
type IrcServer struct {
	servername   string
	serverport   string
	serveraddr   string
	nickname     string
	socket       *net.TCPConn
	protocol     string // can be tcp4, tcp6 , tcp (automatic detect)
	delay        int    // milliseconds to wait when sending Dial commands
	heartbeat    int    // keepalive in seconds
	channel      string // the channel to join
	min_chanlist string // minimal amount of users for a channel to be listed
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
	MyServer.min_chanlist = conf.OConfig["min_chanlist"]

	go MyServer.ircClient()

}

func (this *IrcServer) ircClient() {

	this.ircDial()
	reader := bufio.NewScanner(this.socket)
	var message string = "NOOP" // always better to initialize I/O strings

	for reader.Scan() {

		defer func() {
			if e := recover(); e != nil {
				this.ircDial()
			}
		}()

		message = reader.Text()

		err := reader.Err()

		if err != nil {
			log.Println("[IRC] Error reading socket: %s ", err)
			this.ircDial()
		}

		// does all
		this.IrcInterpreter(message)

		// This must be the last one
		if err == nil {
			log.Printf("[IRC]>  <%s> ", message)
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
		this.IrcCmd("CAP LS")
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		this.IrcCmd("NICK " + this.nickname)
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		userString := "USER " + this.nickname + " " + this.nickname + " " + this.servername + " :" + this.nickname
		log.Print("[AAA] " + userString)
		this.IrcCmd(userString)
		this.IrcCmd("CAP END")

		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		log.Println("[AAA] AAA terminated, now joining")
		this.IrcCmd("JOIN " + this.channel)
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		log.Println("[AAA] Asking for a list of channels")
		this.IrcCmd("LIST >" + this.min_chanlist + ",<10000")

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

func (this *IrcServer) IrcCmd(command string) {

	log.Printf("[IRC] Sending %s command", command)
	_, err := this.socket.Write([]byte(command + "\n"))
	if err != nil {
		log.Printf("[IRC] Cannot send <%s> :  %s:", command, err)

	} else {
		log.Printf("[IRC] Successfully sent <%s>", command)
	}

}

func IrcInitialize() {
	log.Println("[AAA] Initializing IRC Engine...")
}
