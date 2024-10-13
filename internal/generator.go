package internal

import (
	"errors"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Generator struct {
	templates *template.Template
}

func NewGenerator(templates *template.Template) Generator {
	return Generator{
		templates: templates,
	}
}

func (g *Generator) Run() error {
	err := os.Mkdir(buildDirectory, 0755)
	if !errors.Is(err, os.ErrExist) && err != nil {
		return err
	}

	err = os.CopyFS(path.Join(buildDirectory, "static"), os.DirFS(path.Join(contentDirectory, "static")))
	if err != nil {
		return err
	}

	err = filepath.WalkDir(contentDirectory, g.renderTemplates)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) renderTemplates(path string, d fs.DirEntry, err error) error {
	if d.IsDir() {
		return nil
	}

	if !strings.HasSuffix(path, ".html") {
		log.Printf("Skipping %s", path)
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	outputPath := strings.Replace(path, contentDirectory, buildDirectory, 1)
	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if !errors.Is(err, os.ErrExist) && err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = g.templates.ExecuteTemplate(f, "layout.html", PageData{
		Body:        template.HTML(string(data)),
		ProjectPath: "/fullstackgo",
	})
	if err != nil {
		return err
	}

	return nil
}
