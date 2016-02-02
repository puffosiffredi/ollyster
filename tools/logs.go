package tools

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type ollysterlogfile struct {
	logpath  string
	filename string
	logfile  *os.File
}

func init() {

	// just the first time

	//

	var mylogfile ollysterlogfile
	mylogfile.logpath = Hpwd()
	os.MkdirAll(filepath.Join(mylogfile.logpath, "logs"), 0755)

	mylogfile.SetLogFolder()
	go mylogfile.RotateLogFolder()

}

// rotates the log folder

func (this *ollysterlogfile) RotateLogFolder() {

	for {

		time.Sleep(1 * time.Hour)
		if this.logfile != nil {
			err := this.logfile.Close()
			log.Println("[LOG] close logfile returned: ", err)
		}

		this.SetLogFolder()

	}

}

// sets the log folder

func (this *ollysterlogfile) SetLogFolder() {

	const layout = "2006-01-02.15"

	orario := time.Now()

	this.filename = filepath.Join(this.logpath, "logs", "ollyster."+orario.Format(layout)+"00.log")
	log.Println("[LOG] Logfile is: " + this.filename)

	this.logfile, _ = os.Create(this.filename)

	log.SetPrefix("ollyster> ")
	log.SetOutput(this.logfile)

}

// Hpwd: the UNIX pwd
func Hpwd() string {

	tmpLoc, err := os.Getwd()

	if err != nil {
		tmpLoc = "/tmp"
	}

	return tmpLoc

}

func Log_Engine_Start() {

	log.Println("[LOG] LogRotation engine started")

}
