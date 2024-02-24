package templates

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed *
var templateFS embed.FS

func Render(w http.ResponseWriter, t string, data interface{}, renderBase bool) {
	partials := []string{
		"base.layout.html",
		"header.layout.html",
		"footer.layout.html",
		"pie.html",
		"nav.html",
		"expense-stats-grid.html",
	}

	if !renderBase {
		fmt.Println("Rendering without base", t)
		tmpl, err := template.ParseFS(templateFS, t)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

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
