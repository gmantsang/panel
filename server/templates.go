package server

import "github.com/flosch/pongo2"

// Templates holds compiled templates for this app
type Templates struct {
	Index *pongo2.Template
}

// CompileTemplates compiles the templates from source
func CompileTemplates() (*Templates, error) {
	index, err := pongo2.FromFile("templates/index.html")
	if err != nil {
		return nil, err
	}

	return &Templates{
		Index: index,
	}, nil
}
