package main

import (
	"html"
	"log"
	"net/http"
	"os"
)

func main() {
	// Simple static webserver:

	mux := http.NewServeMux()

	mux.HandleFunc("/static/", ServeStatic)
	mux.HandleFunc("/network/", ServeNetwork)
	mux.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe(":8181", mux))
}

// Hpwd: the UNIX pwd
func Hpwd() string {

	tmpLoc, err := os.Getwd()

	if err != nil {
		tmpLoc = "/tmp"
	}

	return tmpLoc

}

//cath-all function for sending people back to homepage.
func home(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/static", http.StatusMovedPermanently)

	log.Printf("Hello, %q\n", html.EscapeString(r.URL.Path))

}

// to manage static contents
func ServeStatic(w http.ResponseWriter, r *http.Request) {

	HttpRoot := Hpwd()
	log.Println("DocumentRoot: ", HttpRoot)
	log.Println("Serving: ", r.URL.Path)
	http.ServeFile(w, r, HttpRoot+r.URL.Path)

}

//to hide logics behind of network.
func ServeNetwork(w http.ResponseWriter, r *http.Request) {

	HttpRoot := Hpwd()
	log.Println("DocumentRoot: ", HttpRoot)
	log.Println("Serving: ", r.URL.Path)
	http.ServeFile(w, r, HttpRoot+"/static/network.html")

}
