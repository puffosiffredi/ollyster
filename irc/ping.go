package irc

import (
	"log"
	"time"
	"ollyster/files"
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
			}
		}()

		time.Sleep(2 * time.Minute)
		log.Printf("[IRC] sending PING :%s" this.servername)

		_, err := this.socket.Write([]byte("PING :" + this.servername + "\n"))
		if err != nil {
			this.socket.Close()
			log.Println("[TCP] Connection down.")
			
			time.Sleep(30 * time.Second)
			log.Println("[TCP] Restarting connection")
			this.ircDial()
		}

	}
}

func (this *IrcServer) ChannelThread() {

	log.Println("[IRC] Initializing the Channel thread")

	for {
		// make it robust

		defer func() {
			if e := recover(); e != nil {
				log.Println("[TCP] Network issue, RECOVER in act")
				
			}
		}()

		time.Sleep(5 * time.Minute)
		files.MyStream.InitializeChanList()
		log.Println("[IRC] Asking for a list of channels")
		this.IrcCmd("LIST >" + this.min_chanlist + ",<10000")
	
	}
}





