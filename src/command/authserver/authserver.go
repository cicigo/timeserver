package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"utils"
	"time"
)

var concurrent_map = utils.NewConcurrentMap()
var config = utils.GetConfig()


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

func dumpContinuously(ticker *time.Ticker, quit chan bool, dumpFile string) {
	for {
		select {
		case <-ticker.C:
			b, error := json.Marshal(concurrent_map.GetData())
			if error != nil {
				log.Errorf("Failed to dump auth info: %s", error)
			}
			dump(b, dumpFile)
		case <- quit:
			ticker.Stop()
			return
		}
		
	}
}

func dump(data []byte, dumpFile string) {
	_, error := os.Stat(dumpFile)
	if error == nil {
		dumpFileBak := dumpFile + ".bak"
		if error = os.Rename(dumpFile, dumpFileBak); error != nil {
			log.Errorf("Back up file failed: %s\n", error)
			return
		}
		if error = ioutil.WriteFile(dumpFile, data, 0644); error != nil {
			log.Errorf("Write file failed: %s\n", error)
			return
		}
		if error = os.Remove(dumpFileBak); error != nil {
			log.Errorf("Delete backup file failed: %s\n", error)
			return
		}
	} else if os.IsNotExist(error) {
		if error = ioutil.WriteFile(dumpFile, data, 0644); error != nil {
			log.Errorf("Write file failed: %s\n", error)
			
		}
	} else {
		log.Errorf("Get file stat failed: %s\n", error)
	}
}

func load(dumpFile string) {
	if _, error := os.Stat(dumpFile); error == nil {
		jsonBlob, err := ioutil.ReadFile(dumpFile)
		if err != nil {
			log.Errorf("Reading auth info file failed, %s.\n", err)
			return
		}
		
		data := make(map[string]string)
		if jsonBlob != nil {
			error = json.Unmarshal(jsonBlob, &data)
			if error == nil {
				concurrent_map.SetData(data)
			} else {
				log.Errorf("Unmarshalling json failed, %s", error)
			}
		}
	}
}

func main() {

	logger, err := log.LoggerFromConfigAsFile(config.Log)

	if err != nil {
		fmt.Printf("configure log file: %s\n", err)
		return
	}

	log.ReplaceLogger(logger)
	
	if config.DumpFile != "" {
		load(config.DumpFile)
	}

	if config.CheckpointInterval > 0 && config.DumpFile != "" {
		log.Infof("Dump auth info every %v seconds.", config.CheckpointInterval)
		ticker := time.NewTicker(time.Duration(config.CheckpointInterval) * time.Second)
		quit := make(chan bool)
		defer func() {quit <- true}()
		go dumpContinuously(ticker, quit, config.DumpFile)
	}

	http.HandleFunc("/", handleCookie)

	error := http.ListenAndServe(fmt.Sprintf(":%v", config.AuthPort), nil)
	if error != nil {
		log.Criticalf("Start auth server with port %d failed: %v\n", config.AuthPort, error)
		os.Exit(1)
	}
	
}
