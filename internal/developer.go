package internal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type DeveloperMode struct {
	mux *http.ServeMux
}

func NewDeveloperMode() DeveloperMode {
	mux := http.NewServeMux()

	staticDir := http.Dir(path.Join(contentDirectory, "static"))
	fi, err := os.Stat(string(staticDir))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fi.IsDir())

	fs := http.FileServer(staticDir)
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /{path...}", servePage)

	return DeveloperMode{
		mux: mux,
	}
}

func renderFromPath(subPath string, w http.ResponseWriter) error {
	indexPath := filepath.Join(contentDirectory, subPath, "index.html")
	log.Printf("Requesting page from '%s'", indexPath)

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	content := strings.ReplaceAll(string(data), "{{projectpath}}", "")

	t, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(w, "layout.html", PageData{
		Body:        template.HTML(content),
		ProjectPath: "",
	})
	if err != nil {
		return err
	}

	return nil
}

func servePage(w http.ResponseWriter, r *http.Request) {
	subPath := r.PathValue("path")
	if subPath == "favicon.ico" {
		w.WriteHeader(200)
		return
	}

	err := renderFromPath(path.Join("/", subPath), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (dm *DeveloperMode) Start() error {
	port := 60600
	log.Printf("Starting developer mode at http://localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), dm.mux)
	if err != nil {
		return err
	}

	return nil
}
