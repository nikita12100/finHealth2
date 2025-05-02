package routes

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

func getTemplate(path string) *template.Template {
	content, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Error trying read file", "path", path, "error", err)
		return nil
	}

	t := template.New(path)
	tmp, err := t.Parse(string(content))
	if err != nil {
		slog.Error("Error parsing html template", "html", string(content), "error", err)
		return nil
	}

	return tmp
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handle /")

	tmp := getTemplate("./static/home_page.html")

	tmp.Execute(w, nil)
}
