package main

import (
	"log"
	"net/http"
	"ollyster/conf"
	"ollyster/files"
	"ollyster/irc"
	"ollyster/tools"
	"ollyster/web"
)

func init() {

	tools.Log_Engine_Start()
	web.InitTmpl()
	files.StreamEngineStart()
	files.InitMsgTmpl()
	conf.StartConfig()
	irc.IrcInitialize()

}

func main() {
	// Simple static webserver:
	log.Println("[WEB] Starting http server on port " + conf.OConfig["webport"])
	mux := http.NewServeMux()

	mux.HandleFunc("/static/", web.ServeStatic)
	mux.HandleFunc("/network/", web.ServeNetwork)
	mux.HandleFunc("/", web.Home)

	log.Fatal(http.ListenAndServe(":"+conf.OConfig["webport"], mux))

}
