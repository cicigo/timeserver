package utils

import (
	log "github.com/cihub/seelog"
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const COOKIE_NAME string = "UUID"

// generate an universally unique identifier
func Uuid() string {
	out, error := exec.Command("/usr/bin/uuidgen").Output()
	if error != nil {
		log.Errorf("Uuid generation failed: %s.", error)
	}
	return strings.Trim(string(out[:]), "\n ")
}

// get UUID from request cookie
func GetUUIDFromCookie(r *http.Request) (string, bool) {
	cookie, error := r.Cookie(COOKIE_NAME)

	if error != nil {
		log.Info("No cookie found")
		return "", false
	} else {
		return cookie.Value, true
	}
}

// get login name from cookie
func GetNameFromCookie(r *http.Request, loggedInNames map[string]string, mutex *sync.Mutex) (string, bool) {
	if uuid, ok := GetUUIDFromCookie(r); ok {
		mutex.Lock()
		defer mutex.Unlock()
		name, nameOk := loggedInNames[uuid]
		return name, nameOk
	} else {
		return "", false
	}

}

func RenderTemplate(w http.ResponseWriter, templatesFolder string, templateName string, data interface{}) {
	tmpl, err := template.New("MyTemplate").ParseFiles(
		filepath.Join(templatesFolder, "framework.html"),
		filepath.Join(templatesFolder, templateName))

	if err != nil {
		log.Criticalf("parsing template files failed: %s\n", err)
	}
	tmpl.ExecuteTemplate(w, "frameworkTemplate", data)
	if err != nil {
		log.Criticalf("executing template failed: %s\n", err)
		return
	}
}
