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
	
	// make it robust

		defer func() {
			if e := recover(); e != nil {
				log.Println("[TCP][PING] Network issue, RECOVER in act")
				time.Sleep(30 * time.Second)
				log.Println("[TCP][PING][REC] Trying to reconnect.")
				this.ircDial()

			}
		}()
	

	log.Println("[IRC][PING] Initializing the KeepAlive engine")

	const layout = "2006-01-02.03:04:05"

	for {
		time.Sleep(2 * time.Minute)
		orario := time.Now()
		log.Printf("[IRC][PING] sending PING :%s", orario.Format(layout))

		_, err := this.socket.Write([]byte("PING :" + orario.Format(layout) + "\n"))

		if err != nil {
			log.Println("[TCP][PING] Network issue, RECOVER in act")
			time.Sleep(10 * time.Second)
			log.Println("[TCP][PING] Trying to reconnect.")
			this.ircDial()

		}

	}
}

func (this *IrcServer) ChannelThread() {
	
	
		// make it robust

		defer func() {
			if e := recover(); e != nil {
				log.Println("[TCP][LIST] Network issue, RECOVER in act")
				
			}
		}()

	log.Println("[IRC][LIST] Initializing the Channel thread")
	time.Sleep(2 * time.Minute)
	for {
	

		files.MyStream.InitializeChanList()
		log.Println("[IRC][LIST] Asking for a list of channels")
		this.IrcCmd("LIST >" + this.min_chanlist + ",<10000")
		time.Sleep(60 * time.Minute)

	}
}
