package irc

import (
	"log"
	"time"
	"ollyster/files"
	"bufio"
)


func init(){
	
	
	go MyServer.KeepAliveThread()
	go MyServer.ChannelThread()

	
}



// rotates the name of streamfiles.
func (this *IrcServer) KeepAliveThread() {

	log.Println("[IRC] Initializing the KeepAlive engine")

	for {
		// make it robust

		defer func() {
			if e := recover(); e != nil {
				log.Println("[TCP] Network issue, RECOVER in act")
				this.socket = nil
				this.socket.Close()
				time.Sleep(30 * time.Second)
				log.Println("[TCP] Trying to reconnect.")
				this.ircDial()
				this.reader = bufio.NewScanner(this.socket)
			}
		}()

		time.Sleep(2 * time.Minute)
		log.Printf("[IRC] sending PING :%s", this.servername)

		 this.socket.Write([]byte("PING :" + this.servername + "\n"))
		

	}
}

func (this *IrcServer) ChannelThread() {

	log.Println("[IRC] Initializing the Channel thread")
    time.Sleep(2 * time.Minute)
	for {
		// make it robust

		defer func() {
			if e := recover(); e != nil {
				log.Println("[TCP] Network issue, RECOVER in act")
				
			}
		}()

		
		files.MyStream.InitializeChanList()
		log.Println("[IRC] Asking for a list of channels")
		this.IrcCmd("LIST >" + this.min_chanlist + ",<10000")
		time.Sleep(60 * time.Minute)
	
	}
}





