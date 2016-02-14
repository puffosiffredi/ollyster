package main

import (
	"html"
	"log"
	"net/http"
	"ollyster/conf"
	"ollyster/files"
	"ollyster/irc"
	"ollyster/tools"
)

func init() {

	tools.Log_Engine_Start()
	files.StreamEngineStart()
	conf.StartConfig()
	irc.IrcInitialize()

}

func main() {
	// Simple static webserver:
	log.Println("[WEB] Starting http server ...")
	mux := http.NewServeMux()

	mux.HandleFunc("/static/", ServeStatic)
	mux.HandleFunc("/network/", ServeNetwork)
	mux.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe(":"+conf.OConfig["webport"], mux))
}

//cath-all function for sending people back to homepage.
func home(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/static", http.StatusMovedPermanently)

	log.Printf("[WEB] Hello, %q\n", html.EscapeString(r.URL.Path))

}

// to manage static contents
func ServeStatic(w http.ResponseWriter, r *http.Request) {

	HttpRoot := tools.Hpwd()
	log.Println("[WEB] DocumentRoot: ", HttpRoot)
	log.Println("[WEB] Serving: ", r.URL.Path)
	http.ServeFile(w, r, HttpRoot+r.URL.Path)

}

//to hide logics behind of network.
func ServeNetwork(w http.ResponseWriter, r *http.Request) {

	HttpRoot := tools.Hpwd()
	log.Println("[WEB] DocumentRoot: ", HttpRoot)
	log.Println("[WEB] Serving: ", r.URL.Path)
	http.ServeFile(w, r, HttpRoot+"/static/network.html")

}
