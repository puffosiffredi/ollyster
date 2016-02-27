package irc

import (
	"bufio"
	"fmt"
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
	MyServer.servername = conf.GetConfItem("servername")
	MyServer.serverport = conf.GetConfItem("serverport")
	MyServer.serveraddr = MyServer.Resolve(MyServer.servername)
	MyServer.nickname = conf.GetConfItem("nickname")
	MyServer.delay, _ = strconv.Atoi(conf.GetConfItem("delay"))
	MyServer.heartbeat, _ = strconv.Atoi(conf.GetConfItem("heartbeat"))
	MyServer.channel = conf.GetConfItem("channel")
	MyServer.min_chanlist = conf.GetConfItem("min_chanlist")

	go MyServer.ircClient()

}

func (this *IrcServer) ircClient() {

	this.ircDial()

	var message string = "NOOP" // always better to initialize I/O strings

	for {

		defer func() {
			if e := recover(); e != nil {
				log.Println("[IRC][REC] Network issue, waiting the network be back")
				err, ok := e.(error)
				if !ok {
					err = fmt.Errorf("[EXC]: %v", e)
				}
				log.Printf("[IRC][REC] Error: <%s>", err)

			}
		}()

		if this.reader.Scan() {
			message = this.reader.Text()
			this.IrcInterpreter(message)
		} else {

			err := this.reader.Err()

			log.Printf("[IRC][READ] Error reading socket: %s ", err)
			log.Println("[IRC][READ] Waiting the connection to be back ")
			time.Sleep(time.Duration(10000 * time.Millisecond))
		}

	}

}

func (this *IrcServer) ircDial() {

	var err error

	defer func() {
		if e := recover(); e != nil {
			log.Println("[TCP][DIAL][REC] Network issue, RECOVER in act")
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("[EXC]: %v", e)
			}
			log.Printf("[TCP][DIAL][REC] Error: <%s>", err)

		}
	}()

	this.socket, err = net.DialTCP(this.protocol, this.ReadIpFromHost(), this.MakeAddr())
	if err != nil {
		log.Printf("[AAA][TCP] CONNECTION ERROR: %s", err)
	} else {
		this.socket.SetKeepAlive(true)                                              // keepalive on
		this.socket.SetKeepAlivePeriod(time.Duration(this.heartbeat) * time.Second) // tcp keepalive to 10 seconds
		this.socket.SetLinger(0)                                                    // brutal close
		this.reader = bufio.NewScanner(this.socket)
		this.writer = bufio.NewWriter(this.socket)

		log.Printf("[AAA][TCP] Connect OK: %s <-> %s", this.socket.LocalAddr().String(), this.socket.RemoteAddr().String())

		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		log.Println("[AAA] Now sending the AAA")
		this.IrcCmd("CAP LS")
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		this.IrcCmd("NICK " + this.nickname)
		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		userString := "USER " + this.nickname + " " + this.nickname + " " + this.servername + " :" + this.nickname
		log.Print("[AAA][USER] " + userString)
		this.IrcCmd(userString)
		this.IrcCmd("CAP END")

		time.Sleep(time.Duration(this.delay) * time.Millisecond)
		log.Println("[AAA] AAA terminated, now joining")
		this.IrcCmd("JOIN " + this.channel)

		time.Sleep(time.Duration(this.delay) * time.Millisecond)

	}

}

func (this *IrcServer) IrcCmd(command string) {

	defer func() {
		if e := recover(); e != nil {
			log.Println("[IRC][CMD][REC] Network issue, cannot write in the socket.")
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("[EXC]: %v", e)
			}
			log.Printf("[IRC][CMD][REC] Error: <%s>", err)

		}
	}()

	log.Printf("[IRC][CMD] Sending %s command", command)

	_, err := this.writer.WriteString(command + "\n")
	if err != nil {
		log.Printf("[IRC][CMD] Cannot send <%s> :  %s:", command, err)
	} else {
		this.writer.Flush()
		log.Printf("[IRC][CMD] Successfully sent <%s>", command)
	}

}

func IrcInitialize() {
	log.Println("[AAA] Initializing IRC Engine...")
}
