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
	reader       *bufio.Scanner
	writer       *bufio.Writer
}

var MyServer IrcServer

func init() {

	// taken by the conf file
	MyServer.servername = conf.OConfig["servername"]
	MyServer.serverport = conf.OConfig["serverport"]
	MyServer.serveraddr = MyServer.Resolve(MyServer.servername)
	MyServer.nickname = conf.OConfig["nickname"]
	MyServer.delay, _ = strconv.Atoi(conf.OConfig["delay"])
	MyServer.heartbeat, _ = strconv.Atoi(conf.OConfig["heartbeat"])
	MyServer.channel = conf.OConfig["channel"]
	MyServer.min_chanlist = conf.OConfig["min_chanlist"]

	go MyServer.ircClient()

}

func (this *IrcServer) ircClient() {

	this.ircDial()

	var message string = "NOOP" // always better to initialize I/O strings

	for this.reader = bufio.NewScanner(this.socket); true; this.reader.Scan() {

		defer func() {
			if e := recover(); e != nil {
				log.Println("[IRC] Network issue, re-dial after 10 sec.")
				this.socket.Close()
				time.Sleep(time.Duration(10000 * time.Millisecond))
				this.ircDial()

			}
		}()

		message = this.reader.Text()

		err := this.reader.Err()
		if err != nil {
			log.Println("[IRC] Error reading socket: %s ", err)
		}

		// does all
		this.IrcInterpreter(message)

	}

}

func (this *IrcServer) ircDial() {

	var err error

	defer func() {
		if e := recover(); e != nil {
			log.Println("[TCP] Network issue, RECOVER in act")
			this.socket = nil
		}
	}()

	this.socket, err = net.DialTCP(this.protocol, nil, this.MakeAddr())

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

	}

}

func (this *IrcServer) IrcCmd(command string) {

	this.writer = bufio.NewWriter(this.socket)

	log.Printf("[IRC] Sending %s command", command)

	_, err := this.writer.WriteString(command + "\n")
	if err != nil {
		log.Printf("[IRC] Cannot send <%s> :  %s:", command, err)
	} else {
		this.writer.Flush()
		log.Printf("[IRC] Successfully sent <%s>", command)
	}

}

func IrcInitialize() {
	log.Println("[AAA] Initializing IRC Engine...")
}
