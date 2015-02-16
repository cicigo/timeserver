package main

import (
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"net/http"
	"os"
	"utils"
)

var concurrent_map = utils.NewConcurrentMap()

func checkRequestParameter(p []string) bool {
	return p != nil && len(p) > 0 && p[0] != ""
}

func handleCookie(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/get" {
		log.Info("Handling GetCookie request.")
		r.ParseForm()
		uuids := r.Form["cookie"]
		if checkRequestParameter(uuids) {
			name := concurrent_map.Get(uuids[0])
			log.Infof("Get name of uuid %s: %s", uuids[0], name)

			fmt.Fprint(w, name)
			w.WriteHeader(200)

		} else {
			log.Warnf("UUID is emtpy")
			w.WriteHeader(400)

		}

	} else if r.URL.Path == "/set" {
		log.Info("Handling PutCookie request.")
		r.ParseForm()
		uuids := r.Form["cookie"]
		names := r.Form["name"]

		if uuids == nil || len(uuids) == 0 {
			log.Warnf("UUID is empty.")
			w.WriteHeader(400)
		} else if names == nil || len(names) == 0 {
			log.Infof("Delete uuid %s.", uuids[0])
			concurrent_map.Delete(uuids[0])
			w.WriteHeader(200)

		} else {
			log.Infof("Put name of uuid %s: %s", uuids[0], names[0])

			concurrent_map.Put(uuids[0], names[0])
			w.WriteHeader(200)
		}

	} else {
		log.Infof("Page Not Found: %s", r.URL.Path)
		w.WriteHeader(404)
	}

}

func main() {
	portPtr := flag.Int("port", 7070, "auth server port number")
	logPtr := flag.String("log", "etc/log.xml", "Log configuration file path")

	flag.Parse()

	logger, err := log.LoggerFromConfigAsFile(*logPtr)

	if err != nil {
		fmt.Printf("configure log file: %s\n", err)
		return
	}

	log.ReplaceLogger(logger)

	http.HandleFunc("/", handleCookie)

	error := http.ListenAndServe(fmt.Sprintf(":%v", *portPtr), nil)
	if error != nil {
		log.Criticalf("Start auth server with port %d failed: %v\n", *portPtr, error)
		os.Exit(1)
	}
}
