package utils

import (
	log "github.com/cihub/seelog"
	"html/template"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
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
func GetUUIDFromCookie(r *http.Request) (string, error) {
	cookie, error := r.Cookie(COOKIE_NAME)

	if error != nil {
		log.Info("No cookie found")
		return "", error
	} else {
		return cookie.Value, nil
	}
}

// get login name from cookie
func GetNameFromCookie(r *http.Request, authClient *AuthClient) (string, error) {
	uuid, error := GetUUIDFromCookie(r)
	if error != nil {
		return "", error
	}

	name, error := authClient.Get(uuid)
	if error != nil {
		return "", error
	}

	return name, nil
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
