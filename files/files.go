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
	streamname string
	streampath string
	streamfile *os.File
}

var MyStream ollysterSocial

func init() {
	go MyStream.RotateSocialFolder()

}

// rotates the name of streamfiles.
func (this *ollysterSocial) RotateSocialFolder() {

	this.streampath = filepath.Join(tools.Hpwd(), "data")
	log.Println("[TXT] Streampath is: " + this.streampath)
	err := os.MkdirAll(filepath.Join(this.streampath), 0755)

	if err != nil {
		log.Printf("[TXT] Cannot create directory: %s", err)

	}

	for {

		const layout = "2006-01-02"
		orario := time.Now()

		this.streamname = filepath.Join(this.streampath, "ollyster."+orario.Format(layout)+".html")
		log.Println("[TXT] Streamfile is: " + this.streamname)
		time.Sleep(10 * time.Minute)

	}

}

// writes down messages for the group
func (this *ollysterSocial) WriteMsgGroup(ev string, gr string, ms string) {

	var err error

	const socialEvent = `
	<li class="list-group-item list-group-item-info">	
	<h4 class="list-group-item-heading">{{.Author}} posted on {{.Group}}</h4>
	<hr>
    <p class="list-group-item-text">{{.Message}}</p>			
	</li> `

	this.streamfile, err = os.OpenFile(this.streamname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("[TXT] Error opening file ", err)
	}

	eventString := string(socialEvent)

	eventString = strings.Replace(eventString, "{{.Author}}", ev, 1)
	eventString = strings.Replace(eventString, "{{.Group}}", gr, 1)
	eventString = strings.Replace(eventString, "{{.Message}}", ms, 1)

	this.AddLineTopFile(eventString)

}

// writes down messages for the group
func (this *ollysterSocial) WriteMsgPriv(ev string, ms string) {

	const socialEvent = `
	<li class="list-group-item list-group-item-success">	
	<h4 class="list-group-item-heading">Private message from <b>{{.Author}}</b></h4>
	<hr>
    <p class="list-group-item-text">{{.Message}}</p>			
	</li> `

	eventString := string(socialEvent)

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

// AddLineToFile : appends one line to the given file.
// only when the line doesn't exists already
// useful when adding groups on the list of groups, or users to the list of users
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

// just starts the engine
func StreamEngineStart() {

	log.Println("[TXT] Stream engine started")

}
