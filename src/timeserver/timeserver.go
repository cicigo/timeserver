//timeserver serves a web page displaying the current time
//of day. The default port number for the webserver is 8080.
//Timeserver only displays time for the time request.
//Using command-line argument -v can show the version
//number.
//
//Copyright 2015 Cici, Chunchao Zhang
package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"utils"
)

var loggedInNames = make(map[string]string)
var mutex = &sync.Mutex{}

const COOKIE_NAME string = "UUID"

type TimeContent struct {
	Time    string
	UtcTime string
	Name    string
}

// set up webpage format and display the current time
func handleTime(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling time request.")

	const layout = "3:04:05PM"
	t := time.Now()
	const utcLayout = "15:04:05 UTC"
	utc := t.UTC()

	name, _ := utils.GetNameFromCookie(r, loggedInNames, mutex)

	timeContent := TimeContent{
		Time:    t.Format(layout),
		UtcTime: utc.Format(utcLayout),
		Name:    name,
	}

	utils.RenderTemplate(w, "templates/time.html", timeContent)
}

//handleNotFound: customarized 404 page for non-time request
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling NotFound URL: %s\n", r.URL.Path)
	w.WriteHeader(404)
	utils.RenderTemplate(w, "templates/notfound.html", nil)
}

// homepage handler
func handleHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		handleNotFound(w, r)
		return
	}

	log.Println("Handling homepage request.")

	if name, ok := utils.GetNameFromCookie(r, loggedInNames, mutex); ok {
		utils.RenderTemplate(w, "templates/greeting.html", name)
	} else {
		utils.RenderTemplate(w, "templates/login.html", nil)
	}
}

// login page handler
func handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling login request.")
	name := html.EscapeString(r.FormValue("name"))
	if name == "" {
		log.Println("log in name is empty")
		utils.RenderTemplate(w, "templates/emptyName.html", nil)
	} else {
		log.Println("log in name is", name)

		uuid := utils.Uuid()
		mutex.Lock()
		loggedInNames[uuid] = name
		mutex.Unlock()
		// Set cookie

		cookie := http.Cookie{Name: COOKIE_NAME, Value: uuid}
		http.SetCookie(w, &cookie)

		// redirect to homepage
		http.Redirect(w, r, "/", 302)
		//fmt.Fprintf(w, fmt.Sprintf("Greeting, %s", name))
	}
}

// logout page handler
func handleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling logout request.")

	// if uuid found in cookie, delete it from loggedInNames
	if uuid, ok := utils.GetUUIDFromCookie(r); ok {
		mutex.Lock()
		delete(loggedInNames, uuid)
		mutex.Unlock()
	}

	// clear cookie
	cookie := http.Cookie{Name: COOKIE_NAME, MaxAge: -1}
	http.SetCookie(w, &cookie)

	// display goodbye message
	utils.RenderTemplate(w, "templates/logout.html", nil)
}

func main() {
	portPtr := flag.Int("port", 8080, "http server port number")
	versionPtr := flag.Bool("v", false, "Display version number")
	flag.Parse()

	if *versionPtr {
		fmt.Println("2.0.0")
		return
	}

	// logger, err := log.LoggerFromConfigAsFile("etc/log.xml")

	// if err != nil {
	// 	fmt.Printf("configure log file: %s\n", err)
	// 	return
	// }

	// log.ReplaceLogger(logger)

	// handle css
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.HandleFunc("/time", handleTime)
	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/index.html", handleHomePage)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)

	error := http.ListenAndServe(fmt.Sprintf(":%v", *portPtr), nil)
	if error != nil {
		fmt.Printf("Start server with port %d failed: %v\n", *portPtr, error)
		os.Exit(1)
	}
}
