package radios

import "github.com/flosch/pongo2"

// Templates contains the templates for radios
type Templates struct {
	List *pongo2.Template
	Edit *pongo2.Template
}

// CompileTemplates builds the pongo2 templates for radios
func CompileTemplates() (*Templates, error) {
	list, err := pongo2.FromFile("templates/radio/list.html")
	if err != nil {
		return nil, err
	}

	edit, err := pongo2.FromFile("templates/radio/edit.html")
	if err != nil {
		return nil, err
	}

	return &Templates{
		List: list,
		Edit: edit,
	}, nil
}
