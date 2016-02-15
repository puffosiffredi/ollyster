package web

import (
	"io"
	"log"
	"net/http"
	fi "ollyster/files"
	"ollyster/tools"
	"strings"
)

// to manage static contents
func ServeStatic(w http.ResponseWriter, r *http.Request) {

	HttpRoot := tools.Hpwd()
	log.Println("[WEB] DocumentRoot: ", HttpRoot)
	log.Println("[WEB] Static Serving: ", r.URL.Path)

	if (r.URL.Path == "/static/") || (r.URL.Path == "/static") {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		log.Println("[WEB] No you don't")
	} else {

		http.ServeFile(w, r, HttpRoot+r.URL.Path)
	}
}

// to manage the home page
// 	io.WriteString(w, st)
func Home(w http.ResponseWriter, r *http.Request) {

	HttpRoot := tools.Hpwd()
	log.Println("[WEB] DocumentRoot: ", HttpRoot)
	log.Println("[WEB] Home Serving: ", r.URL.Path)
	// http.ServeFile(w, r, HttpRoot+r.URL.Path)
	
	contents := fi.MyStream.RetrieveStreamString()
	
	pageString := strings.Replace(OTemplates.indextmpl, "{{.Contents}}", contents, 1)
	io.WriteString(w, pageString)

}

//to hide logics behind of network.
func ServeNetwork(w http.ResponseWriter, r *http.Request) {

	HttpRoot := tools.Hpwd()
	log.Println("[WEB] DocumentRoot: ", HttpRoot)
	log.Println("[WEB] Serving: ", r.URL.Path)
	http.ServeFile(w, r, HttpRoot+"/static/network.html")

}
