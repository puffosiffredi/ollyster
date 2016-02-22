package irc

import (

	"log"
	"ollyster/files"
	"time"
)

func init() {

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
				
				this.socket.Close()
				
				time.Sleep(30 * time.Second)
				log.Println("[TCP] Trying to reconnect.")
				this.ircDial()

			}
		}()

		time.Sleep(2 * time.Minute)
		log.Printf("[IRC] sending PING :%s", this.servername)

		_, err := this.socket.Write([]byte("PING :" + this.servername + "\n"))

		if err != nil {
			log.Println("[TCP][PING] Network issue, RECOVER in act")
			
			this.socket.Close()
			
			
			time.Sleep(10 * time.Second)
			log.Println("[TCP] Trying to reconnect.")
			this.ircDial()

		}

	}
}

func (this *IrcServer) ChannelThread() {

	log.Println("[IRC] Initializing the Channel thread")
	time.Sleep(2 * time.Minute)
	for {
		// make it robust

		defer func() {
			if e := recover(); e != nil {
				log.Println("[TCP][LIST] Network issue, RECOVER in act")
			
			this.socket.Close()
			time.Sleep(10 * time.Second)
			log.Println("[TCP][LIST] Trying to reconnect.")
			this.ircDial()

			}
		}()

		files.MyStream.InitializeChanList()
		log.Println("[IRC] Asking for a list of channels")
		this.IrcCmd("LIST >" + this.min_chanlist + ",<10000")
		time.Sleep(60 * time.Minute)

	}
}
