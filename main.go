package main

import (
	"fmt"
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

	defer func() {
		if e := recover(); e != nil {
			log.Println("[MAIN] OH, SHIT.")
			err, ok := e.(error)
			if !ok {
				err = fmt.Errorf("[EXC]: %v", e)
			}
			log.Printf("[MAIN] Error: <%s>", err)

		}
	}()

	// Simple static webserver:
	log.Println("[WEB] Starting http server on port " + conf.GetConfItem("webport"))
	mux := http.NewServeMux()

	mux.HandleFunc("/static/", web.ServeStatic)
	mux.HandleFunc("/network/", web.ServeNetwork)
	mux.HandleFunc("/inbox/", web.ServeInbox)
	mux.HandleFunc("/get/addgroup/", web.AddGroup)
	mux.HandleFunc("/", web.Home)

	log.Fatal(http.ListenAndServe(":"+conf.GetConfItem("webport"), mux))

}
