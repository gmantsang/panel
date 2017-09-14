package server

import "github.com/flosch/pongo2"

// Templates holds compiled templates for this app
type Templates struct {
	Index     *pongo2.Template
	RadioList *pongo2.Template
	RadioEdit *pongo2.Template
}

// CompileTemplates compiles the templates from source
func CompileTemplates() (*Templates, error) {
	index, err := pongo2.FromFile("templates/index.html")
	if err != nil {
		return nil, err
	}
	radioList, err := pongo2.FromFile("templates/radio-list.html")
	if err != nil {
		return nil, err
	}
	radioEdit, err := pongo2.FromFile("templates/radio-edit.html")
	if err != nil {
		return nil, err
	}

	return &Templates{
		Index:     index,
		RadioList: radioList,
		RadioEdit: radioEdit,
	}, nil
}
