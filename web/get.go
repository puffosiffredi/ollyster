package web

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

// to manage GET requests
func AddGroup(w http.ResponseWriter, r *http.Request) {

	if Gandalf(w, r) {

		log.Println("[WEB][GET] Request: ", r.URL.Path)

		re, _ := regexp.Compile("^/get/addgroup/([^/]+)$")
		match := re.FindStringSubmatch(r.URL.Path)

		if match != nil {
			log.Println("[WEB][GET] Requested group: ", strings.Replace(match[1], "@", "#", -1))
			http.Redirect(w, r, "/network/", http.StatusMovedPermanently)

		} else {

			log.Println("[WEB][GET] No requested group!!! ")
			http.Redirect(w, r, "/network/", http.StatusMovedPermanently)

		}

	}

}
