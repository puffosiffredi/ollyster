package web

import (
	"log"
	"net/http"
	co "ollyster/conf"
)

func Gandalf(w http.ResponseWriter, r *http.Request) {

	pass := false

	// WWW-Authenticate: Basic realm="Ollyster"

	username, password, _ := r.BasicAuth()

	log.Printf("[WEB][AUTH] username/password in http : %s/%s", username, password)
	log.Printf("[WEB][AUTH] username/password expected: %s/%s", co.OConfig["username"], co.OConfig["password"])

	pass = (username == co.OConfig["username"] && password == co.OConfig["password"])

	if pass == false {
		log.Println("[WEB][AUTH] YOU SHALL NOT PASS!")
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Ollyster-\"")
		http.Error(w, "authorization failed", http.StatusUnauthorized)
	}

}
