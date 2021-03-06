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

	streamname string // name of the file containing the stream in html: messages, etc
	streampath string // path of the file which contains the stream

	privtname string // name of the file where we store private messages and notices.

	//  file needed for the channel list

	channelname string // name of the file containing a dump of channelnames. Mostly for debug.

	OChannels map[string]string // map containing a list of channelname, channeldesscription

	NamesBuf string // contains all the people which subscribed channels we subscribed

}

var MyStream ollysterSocial

func init() {

	MyStream.InitializeChanList()
	MyStream.streampath = filepath.Join(tools.Hpwd(), "data")

	log.Println("[TXT][INI] Streampath is: " + MyStream.streampath)
	err := os.MkdirAll(filepath.Join(MyStream.streampath), 0755)
	if err != nil {
		log.Printf("[TXT][INI] Cannot create directory: %s", err)
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
		this.privtname = filepath.Join(this.streampath, "inbox."+orario.Format(layout)+".html")
		this.channelname = filepath.Join(this.streampath, "channels.txt")

		log.Println("[TXT] Streamfile is now: " + this.streamname)
		log.Println("[TXT] Channelfile is now: " + this.channelname)
		log.Println("[TXT] Privtfile is now: " + this.privtname)

		// initializes the streamname if it doesn't exists

		_, err := os.Stat(this.streamname)
		if err != nil {
			ioutil.WriteFile(this.streamname, []byte("<!---Rotation Engine was here -->"), 0755)

		}

		// initializes channelname if it doesn't exists

		_, err = os.Stat(this.channelname)
		if err != nil {
			ioutil.WriteFile(this.channelname, []byte("<!---Rotation Engine was here -->"), 0755)

		}

		// initializes privtname if it doesn't exists

		_, err = os.Stat(this.privtname)
		if err != nil {
			ioutil.WriteFile(this.privtname, []byte("<!---Rotation Engine was here -->"), 0755)

		}

		time.Sleep(10 * time.Minute)

	}

}

// writes down messages for the group
func (this *ollysterSocial) WriteMsgGroup(ev string, gr string, ms string) {
	const layout = "Jan 02 15:04"
	orario := time.Now()

	eventString := MyOllysterMsg.AzureTmpl

	eventString = strings.Replace(eventString, "{{.Author}}", ev, -1)
	eventString = strings.Replace(eventString, "{{.Group}}", gr, -1)
	eventString = strings.Replace(eventString, "{{.Message}}", ms, -1)
	eventString = strings.Replace(eventString, "{{.Time}}", orario.Format(layout), -1)

	this.AddLineTopFile(eventString)

}

// writes down messages for the group
func (this *ollysterSocial) WriteMsgMention(ev string, gr string, ms string) {

	const layout = "Jan 02 15:04"
	orario := time.Now()

	eventString := MyOllysterMsg.RedTmpl

	eventString = strings.Replace(eventString, "{{.Author}}", ev, -1)
	eventString = strings.Replace(eventString, "{{.Group}}", gr, -1)
	eventString = strings.Replace(eventString, "{{.Message}}", ms, -1)
	eventString = strings.Replace(eventString, "{{.Time}}", orario.Format(layout), -1)

	this.AddLineTopFile(eventString)

}

// writes down messages for the group
func (this *ollysterSocial) WriteMsgPriv(ev string, ms string) {

	const layout = "Jan 02 15:04"
	orario := time.Now()

	eventString := MyOllysterMsg.GreenTmpl

	eventString = strings.Replace(eventString, "{{.Author}}", ev, -1)
	eventString = strings.Replace(eventString, "{{.Message}}", ms, -1)
	eventString = strings.Replace(eventString, "{{.Time}}", orario.Format(layout), -1)

	// let's keep it private

	this.AddPrivTopFile(eventString)

}

func (this *ollysterSocial) WriteNotice(sender string, msg string) {

	const layout = "Jan 02 15:04"
	orario := time.Now()

	eventString := MyOllysterMsg.YellowTmpl

	eventString = strings.Replace(eventString, "{{.Author}}", sender, -1)
	eventString = strings.Replace(eventString, "{{.Message}}", msg, -1)
	eventString = strings.Replace(eventString, "{{.Time}}", orario.Format(layout), -1)

	// let's keep it private

	this.AddPrivTopFile(eventString)

}

// RetrieveStringFromFile returns a file into a single string
// useful to retrieve the content and shoot into the home page
func (this *ollysterSocial) RetrieveStreamString(sname string) string {

	var myfile string

	if sname == "private" {

		myfile = this.privtname

	} else {

		myfile = this.streamname
	}

	content, err := ioutil.ReadFile(myfile)
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

// Add PrivtTopFile
// AddLineTopFile : appends one line to the privt file, in reversed order, last one top
func (this *ollysterSocial) AddPrivTopFile(line string) error {

	content, err := ioutil.ReadFile(this.privtname)
	if err != nil {
		return err
	}

	contentString := line + "\n" + string(content)

	err = ioutil.WriteFile(this.privtname, []byte(contentString), 0755)
	if err != nil {
		return err
	}

	return nil

}

// adds one user to the list of available users
func (this *ollysterSocial) AddUniqueUser(userline string) error {

	this.NamesBuf += "\n" + userline

	return nil

}

// just starts the engine
func StreamEngineStart() {

	log.Println("[TXT] Stream engine started")

}
