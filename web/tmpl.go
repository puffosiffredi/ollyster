package web

import (
	"io/ioutil"
	"log"
	to "ollyster/tools"
	fp "path/filepath"
	"strings"
)

type ollysterTmpl struct {
	indextmpl string
	grouptmpl string
	profifmpl string
	inboxtmpl string
}

var OTemplates ollysterTmpl

func init() {

	indextmplFile := fp.Join(to.Hpwd(), "static", "tmpl", "index.tmpl")
	grouptmplFile := fp.Join(to.Hpwd(), "static", "tmpl", "group.tmpl")
	profitmplFile := fp.Join(to.Hpwd(), "static", "tmpl", "profile.tmpl")
	inboxtmplFile := fp.Join(to.Hpwd(), "static", "tmpl", "inbox.tmpl")

	OTemplates.indextmpl = RetrieveTmplString(indextmplFile)
	OTemplates.grouptmpl = RetrieveTmplString(grouptmplFile)
	OTemplates.profifmpl = RetrieveTmplString(profitmplFile)
	OTemplates.inboxtmpl = RetrieveTmplString(inboxtmplFile)

}

// useful to retrieve the content and shoot into the home page
func RetrieveTmplString(file string) string {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("[WEB] Cannot retrieve template " + file)
		return "<!-- EMPTY FILE -->"
	}

	log.Println("[WEB] Template " + file + " loaded")
	return string(content)

}

func MakeChannels(channels string) (badges string) {

	sepa := "</span><span class=\"badge\">"
	chans := "<span class=\"badge\">" + strings.Replace(channels, ",", sepa, -1) + "</span>"

	return chans

}

func InitTmpl() {

	log.Println("[WEB] Initializing template engine")

}
