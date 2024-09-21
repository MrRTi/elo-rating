package render

import (
	"html/template"
	"net/http"
)

func RenderTemplate(templatePath string, w http.ResponseWriter, args any) error {
	tmpl := template.Must(template.ParseFiles(templatePath))
	return tmpl.Execute(w, args)
}

func RenderRedirect(url string, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func RenderBadRequest(errorMessage string, w http.ResponseWriter) {
	http.Error(w, errorMessage, http.StatusBadRequest)
}

func RenderMethodNotAllowed(errorMessage string, w http.ResponseWriter) {
	http.Error(w, errorMessage, http.StatusMethodNotAllowed)
}

func RenderNotFound(errorMessage string, w http.ResponseWriter) {
	http.Error(w, errorMessage, http.StatusNotFound)
}
