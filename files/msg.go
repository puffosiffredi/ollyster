package files

import (
	"io/ioutil"
	"log"
	to "ollyster/tools"
	fp "path/filepath"
)

type ollysterMsg struct {
	RedTmpl    string
	GreenTmpl  string
	YellowTmpl string
	AzureTmpl  string
}

var MyOllysterMsg ollysterMsg

func init() {

	redFile := fp.Join(to.Hpwd(), "static", "tmpl", "msg-red.tmpl")
	yellowFile := fp.Join(to.Hpwd(), "static", "tmpl", "msg-yellow.tmpl")
	greenFile := fp.Join(to.Hpwd(), "static", "tmpl", "msg-green.tmpl")
	azureFile := fp.Join(to.Hpwd(), "static", "tmpl", "msg-azure.tmpl")

	MyOllysterMsg.RedTmpl = RetrieveTmplString(redFile)
	MyOllysterMsg.GreenTmpl = RetrieveTmplString(greenFile)
	MyOllysterMsg.YellowTmpl = RetrieveTmplString(yellowFile)
	MyOllysterMsg.AzureTmpl = RetrieveTmplString(azureFile)

}

// useful to retrieve the content and shoot into the home page
func RetrieveTmplString(file string) string {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("[TXT] Cannot retrieve template " + file)
		return "<!-- EMPTY FILE -->"
	}

	log.Println("[TXT] Template " + file + " loaded")
	return string(content)

}

func InitMsgTmpl() {

	log.Println("[TXT] Initializing Message Template engine")

}
