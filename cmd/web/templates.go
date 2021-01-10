package main

import (
	"path/filepath"
	"text/template"

	"orinayooyelade.com/snippetbox/pkg/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		templateSet, err = template.ParseFiles(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		templateSet, err = template.ParseFiles(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
