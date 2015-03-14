//timeserver serves a web page displaying the current time
//of day. The default port number for the webserver is 8080.
//Timeserver only displays time for the time request.
//Using command-line argument -v can show the version
//number.
//
//Copyright 2015 Cici, Chunchao Zhang
package main

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"html"
	"math/rand"
	"net/http"
	"os"
	"time"
	"utils"
	"utils/serverconfig"
)

var config = serverconfig.GetConfig()
var authClient = utils.NewAuthClient(fmt.Sprintf("%s:%v", config.AuthHost, config.AuthPort))
var limiter = utils.NewLimiter(config.MaxInflight)

var loginCounter = utils.NewCounter("login")
var timeUserCounter = utils.NewCounter("time-user")
var timeAnonCounter = utils.NewCounter("time-anon")
var time200sCounter = utils.NewCounter("200s")
var time404sCounter = utils.NewCounter("404s")

const COOKIE_NAME string = "UUID"

type TimeContent struct {
	Time    string
	UtcTime string
	Name    string
}

// load simulator
func load() {
	load := time.Duration((rand.NormFloat64()*config.DeviationMs + config.AvgResponseMs)) * time.Millisecond
	time.Sleep(load)
}

// set up webpage format and display the current time
func handleTime(w http.ResponseWriter, r *http.Request) {
	if !limiter.Get() {
		log.Warn("Cannot get token")
		w.WriteHeader(500)
		return
	}
	defer limiter.Release()

	log.Info("Handling time request.")

	load()
	const layout = "3:04:05PM"
	t := time.Now()
	const utcLayout = "15:04:05 UTC"
	utc := t.UTC()

	name, error := utils.GetNameFromCookie(r, authClient)

	if error != nil {
		// clear cookie
		cookie := http.Cookie{Name: COOKIE_NAME, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

	timeContent := TimeContent{
		Time:    t.Format(layout),
		UtcTime: utc.Format(utcLayout),
		Name:    name,
	}

	if name != "" {
		timeUserCounter.Incr(1)
	} else {
		timeAnonCounter.Incr(1)
	}
	time200sCounter.Incr(1)

	utils.RenderTemplate(w, config.Templates, "time.html", timeContent)
}

//handleNotFound: customarized 404 page for non-time request
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Infof("Handling NotFound URL: %s\n", r.URL.Path)
	time404sCounter.Incr(1)

	w.WriteHeader(404)
	utils.RenderTemplate(w, config.Templates, "notfound.html", nil)
}

// homepage handler
func handleHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		handleNotFound(w, r)
		return
	}

	if !limiter.Get() {
		w.WriteHeader(500)
		return
	}
	defer limiter.Release()

	log.Infof("Handling homepage request.")

	name, error := utils.GetNameFromCookie(r, authClient)

	if error != nil {
		// clear cookie
		cookie := http.Cookie{Name: COOKIE_NAME, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

	time200sCounter.Incr(1)
	if name == "" {
		utils.RenderTemplate(w, config.Templates, "login.html", nil)
	} else {
		utils.RenderTemplate(w, config.Templates, "greeting.html", name)
	}
}

// login page handler
func handleLogin(w http.ResponseWriter, r *http.Request) {
	loginCounter.Incr(1)

	if !limiter.Get() {
		w.WriteHeader(500)
		return
	}
	defer limiter.Release()

	log.Infof("Handling login request.")
	name := html.EscapeString(r.FormValue("name"))
	if name == "" {
		log.Info("log in name is empty")
		time200sCounter.Incr(1)
		utils.RenderTemplate(w, config.Templates, "emptyname.html", nil)
	} else {
		log.Infof("log in name is %s.", name)

		uuid := utils.Uuid()

		if err := authClient.Set(uuid, name); err != nil {
			log.Errorf("log in failed.: %s", err)
			w.WriteHeader(500)
			return
		}

		// Set cookie

		cookie := http.Cookie{Name: COOKIE_NAME, Value: uuid}
		http.SetCookie(w, &cookie)

		// redirect to homepage
		http.Redirect(w, r, "/", 302)
	}
}

// logout page handler
func handleLogout(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling logout request.")

	// if uuid found in cookie, delete it from loggedInNames
	uuid, error := utils.GetUUIDFromCookie(r)

	if error == nil && uuid != "" {
		authClient.Delete(uuid)
	}

	// clear cookie
	cookie := http.Cookie{Name: COOKIE_NAME, MaxAge: -1}
	http.SetCookie(w, &cookie)

	time200sCounter.Incr(1)
	// display goodbye message
	utils.RenderTemplate(w, config.Templates, "logout.html", nil)
}

func handleMonitor(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling monitor request.")

	counters := utils.DumpCounter()
	b, error := json.Marshal(counters)
	if error != nil {
		log.Errorf("Failed to dump auth info: %s", error)
		w.WriteHeader(500)
		return
	}
	fmt.Fprint(w, string(b[:]))
}

func main() {
	if config.Version {
		fmt.Println("2.0.0")
		return
	}

	logger, err := log.LoggerFromConfigAsFile(config.Log)

	if err != nil {
		fmt.Printf("configure log file: %s\n", err)
		return
	}

	log.ReplaceLogger(logger)

	log.Infof("Templates folder is %s.", config.Templates)

	// handle css
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.HandleFunc("/time", handleTime)
	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/index.html", handleHomePage)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/monitor", handleMonitor)

	error := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
	if error != nil {
		log.Criticalf("Start server with port %d failed: %v\n", config.Port, error)
		os.Exit(1)
	}
}
