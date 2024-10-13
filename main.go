package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template

type PageData struct {
	Body template.HTML
}

func init() {
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	templates = t
}

func main() {
	err := filepath.WalkDir("content", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fmt.Println(path, d.Name(), "directory?", d.IsDir())

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		outputPath := strings.Replace(path, "content", "docs", 1)
		err = os.MkdirAll(filepath.Dir(outputPath), 0755)
		if errors.Is(err, os.ErrExist) {

		} else if err != nil {
			return err
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		err = templates.ExecuteTemplate(f, "layout.html", PageData{
			Body: template.HTML(string(data)),
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
