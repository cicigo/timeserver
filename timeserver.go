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
	"net/http"
	"os"
	"time"
	"os/exec"
	"log"
	"sync"
)

var loggedInNames = make(map[string]string)
var mutex = &sync.Mutex{}

func uuid() string {
	out, error := exec.Command("/usr/bin/uuidgen").Output()
	if error != nil {
		log.Fatal(error)
	}
	return string(out[:])
}


//handleTime: set up webpage format and display the current time
func handleTime(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling time request.")
	const layout = "3:04:05PM"
	t := time.Now()
	content := fmt.Sprintf(`
        <html>
	<head>
	<style>
	p {font-size: xx-large}
	span.time {color: red}
	</style>
	</head>
	<body>
	<p>The time is now <span class="time">%s</span>.</p>
	</body>
	</html>`, t.Format(layout))

	fmt.Fprintf(w, content)
}

//handleNotFound: customarized 404 page for non-time request
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling NotFound page.")
	content :=
		`
<html>
<body>
<p>These are not the URLs you're looking for.</p>
</body>
</html>
`
	fmt.Fprintf(w, content)
}

func handleUUID(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling UUID request.")
	fmt.Fprintf(w, uuid())
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Index page request.")
	loginForm := `
<html>
<body>
<form action="login">
  What is your name, Earthling?
  <input type="text" name="name" size="50">
  <input type="submit">
</form>
</p>
</body>
</html>
`
	greetings := `Greetings, %s`
	
	cookie, error := r.Cookie("uid")

	if error != nil {
		log.Println("No cookie found")
		fmt.Fprintf(w, loginForm)
	} else {
		mutex.Lock()
		name, ok := loggedInNames[cookie.Value]
		mutex.Unlock()
		if ok { // name is logged in
			fmt.Fprintf(w, fmt.Sprintf(greetings, name))
		} else {
			fmt.Fprintf(w, loginForm)
		}
	}
	
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling login page.")
	name := r.FormValue("name")
	if name == "" {
		log.Println("log in name is empty")
		fmt.Fprintf(w, "C'mon, I need a name")
	} else {
		log.Println("log in name is", name)
		
		mutex.Lock()
		loggedInNames[name] = uuid()
		mutex.Unlock()
		// TODO: redirect
		fmt.Fprintf(w, fmt.Sprintf("Greeting, %s", name))
	}
	
}

func main() {
	portPtr := flag.Int("port", 8080, "http server port number")
	versionPtr := flag.Bool("v", false, "Display version number")
	flag.Parse()

	if *versionPtr {
		fmt.Println("1.0.0")
		return
	}

	http.HandleFunc("/time", handleTime)
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/index.html", handleIndex)
	http.HandleFunc("/uuid", handleUUID)
	http.HandleFunc("/login", handleLogin)

	error := http.ListenAndServe(fmt.Sprintf(":%v", *portPtr), nil)
	if error != nil {
		fmt.Printf("Start server with port %d failed: %v\n", *portPtr, error)
		os.Exit(1)
	}
}
