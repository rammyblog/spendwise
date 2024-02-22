package templates

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed *
var templateFS embed.FS

func Render(w http.ResponseWriter, t string, data interface{}) {
	fmt.Println("Rendering template: ", t)
	fmt.Println("Data: ", data)
	partials := []string{
		"base.layout.html",
		"header.layout.html",
		"footer.layout.html",
		// "error.html",
		"nav.html",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, t)

	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
