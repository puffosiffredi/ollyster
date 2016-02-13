package main

import (
	"html"
	"log"
	"net/http"
	"ollyster/irc"
	"ollyster/tools"
)

func init() {

	tools.Log_Engine_Start()
	irc.IrcInitialize()

}

func main() {
	// Simple static webserver:
	log.Println("[WEB] Starting http server on port 8181...")
	mux := http.NewServeMux()

	mux.HandleFunc("/static/", ServeStatic)
	mux.HandleFunc("/network/", ServeNetwork)
	mux.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe(":8181", mux))
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
