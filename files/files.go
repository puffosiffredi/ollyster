package files

import (
	"io/ioutil"
	"log"
	"ollyster/tools"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ollysterSocial struct {

	// file needed for the stream

	streamname string
	streampath string
	streamfile *os.File

	//  file needed for the channel list

	channelname string
	Channelbuf  string
}

var MyStream ollysterSocial

func init() {

	MyStream.Channelbuf = "<!--Placeholder for the list of groups-->"
	MyStream.streampath = filepath.Join(tools.Hpwd(), "data")
	log.Println("[TXT] Streampath is: " + MyStream.streampath)
	err := os.MkdirAll(filepath.Join(MyStream.streampath), 0755)
	if err != nil {
		log.Printf("[TXT] Cannot create directory: %s", err)
	} else {
		go MyStream.RotateSocialFolder()
	}

}

// rotates the name of streamfiles.
func (this *ollysterSocial) RotateSocialFolder() {

	for {

		const layout = "2006-01-02"
		orario := time.Now()

		this.streamname = filepath.Join(this.streampath, "ollyster."+orario.Format(layout)+".html")
		this.channelname = filepath.Join(this.streampath, "channels.txt")

		log.Println("[TXT] Streamfile is now: " + this.streamname)
		log.Println("[TXT] Channelfile is now: " + this.channelname)

		// initializes the streamname if it doesn't exists

		_, err := os.Stat(this.streamname)
		if err != nil {
			ioutil.WriteFile(this.streamname, []byte("<!---Rotation Engine was here -->"), 0755)

		}

		_, err = os.Stat(this.channelname)
		if err != nil {
			ioutil.WriteFile(this.channelname, []byte("<!---Rotation Engine was here -->"), 0755)

		}

		time.Sleep(10 * time.Minute)

	}

}

func (this *ollysterSocial) FlushChanList() {

	// initialize the group file if it doesn't exists
	// periodically flushes the channelbuf there

	ioutil.WriteFile(this.channelname, []byte(this.Channelbuf), 0755)

}

// writes down messages for the group
func (this *ollysterSocial) WriteMsgGroup(ev string, gr string, ms string) {

	eventString := MyOllysterMsg.AzureTmpl

	eventString = strings.Replace(eventString, "{{.Author}}", ev, 1)
	eventString = strings.Replace(eventString, "{{.Group}}", gr, 1)
	eventString = strings.Replace(eventString, "{{.Message}}", ms, 1)

	this.AddLineTopFile(eventString)

}

// writes down messages for the group
func (this *ollysterSocial) WriteMsgPriv(ev string, ms string) {

	eventString := MyOllysterMsg.GreenTmpl

	eventString = strings.Replace(eventString, "{{.Author}}", ev, 1)
	eventString = strings.Replace(eventString, "{{.Message}}", ms, 1)

	this.AddLineTopFile(eventString)

}

// RetrieveStringFromFile returns a file into a single string
// useful to retrieve the content and shoot into the home page
func (this *ollysterSocial) RetrieveStreamString() string {

	content, err := ioutil.ReadFile(this.streamname)
	if err != nil {
		return "<!-- EMPTY FILE -->"
	}

	return string(content)

}

// AddLineTopFile : appends one line to the given file, in reversed order, last one top
func (this *ollysterSocial) AddLineTopFile(line string) error {

	content, err := ioutil.ReadFile(this.streamname)
	if err != nil {
		return err
	}

	contentString := line + "\n" + string(content)

	err = ioutil.WriteFile(this.streamname, []byte(contentString), 0755)
	if err != nil {
		return err
	}

	return nil

}

// AddUniqueChannel adds a channel to the channel buffer if it doesn't exists.

func (this *ollysterSocial) AddUniqueChannel(channelline string) error {

	this.Channelbuf += "\n" + channelline

	return nil

}

// just starts the engine
func StreamEngineStart() {

	log.Println("[TXT] Stream engine started")

}
