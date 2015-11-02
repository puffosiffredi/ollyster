package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
)

func main() {
	// Simple static webserver:

	mux := http.NewServeMux()

	mux.HandleFunc("/static/", ServeStatic)
	mux.HandleFunc("/", hello)

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

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
	fmt.Fprintf(w, "If you see this page this is because URL %q is wrong\n", html.EscapeString(r.URL.Path))
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {

	HttpRoot := Hpwd()
	log.Println("DocumentRoot: ", HttpRoot)
	log.Println("Serving: ", r.URL.Path)
	http.ServeFile(w, r, HttpRoot+r.URL.Path)

}
