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
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var loggedInNames = make(map[string]string)
var mutex = &sync.Mutex{}

const COOKIE_NAME string = "UUID"

// generate an universally unique identifier
func uuid() string {
	out, error := exec.Command("/usr/bin/uuidgen").Output()
	if error != nil {
		log.Fatal(error)
	}
	return strings.Trim(string(out[:]), "\n ")
}

// get UUID from request cookie
func getUUIDFromCookie(r *http.Request) (string, bool) {
	cookie, error := r.Cookie(COOKIE_NAME)

	if error != nil {
		log.Println("No cookie found")
		return "", false
	} else {
		return cookie.Value, true
	}
}

// get login name from cookie
func getNameFromCookie(r *http.Request) (string, bool) {
	if uuid, ok := getUUIDFromCookie(r); ok {
		mutex.Lock()
		defer mutex.Unlock()
		name, nameOk := loggedInNames[uuid]
		return name, nameOk
	} else {
		return "", false
	}

}

type TimeContent struct {
	Time string
	Name string
}

// set up webpage format and display the current time
func handleTime(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling time request.")

	const layout = "3:04:05PM"
	t := time.Now()
	name, _ := getNameFromCookie(r)

	timeContent := TimeContent{
		Time: t.Format(layout),
		Name: name,
	}

	renderTemplate(w, "templates/time.html", timeContent)
}

func renderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.New("MyTemplate").ParseFiles("templates/framework.html", templatePath)
	if err != nil {
		fmt.Printf("parsing template files failed: %s\n", err)
	}
	tmpl.ExecuteTemplate(w, "frameworkTemplate", data)
	if err != nil {
		fmt.Printf("executing template failed: %s\n", err)
		return
	}
}

//handleNotFound: customarized 404 page for non-time request
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling NotFound URL: %s\n", r.URL.Path)
	w.WriteHeader(404)
	renderTemplate(w, "templates/notfound.html", nil)
}

// uuid page handler, for testing.
func handleUUID(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling UUID request.")
	fmt.Fprintf(w, uuid())
}

// homepage handler
func handleHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		handleNotFound(w, r)
		return
	}

	log.Println("Handling homepage request.")

	if name, ok := getNameFromCookie(r); ok {
		renderTemplate(w, "templates/greeting.html", name)
	} else {
		renderTemplate(w, "templates/login.html", nil)
	}
}

// login page handler
func handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling login request.")
	name := html.EscapeString(r.FormValue("name"))
	if name == "" {
		log.Println("log in name is empty")
		fmt.Fprintf(w, "C'mon, I need a name")
	} else {
		log.Println("log in name is", name)

		uuid := uuid()
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
	if uuid, ok := getUUIDFromCookie(r); ok {
		mutex.Lock()
		delete(loggedInNames, uuid)
		mutex.Unlock()
	}

	// clear cookie
	cookie := http.Cookie{Name: COOKIE_NAME, MaxAge: -1}
	http.SetCookie(w, &cookie)

	// display goodbye message
	renderTemplate(w, "templates/logout.html", nil)
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
	http.HandleFunc("/uuid", handleUUID)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)

	error := http.ListenAndServe(fmt.Sprintf(":%v", *portPtr), nil)
	if error != nil {
		fmt.Printf("Start server with port %d failed: %v\n", *portPtr, error)
		os.Exit(1)
	}
}
