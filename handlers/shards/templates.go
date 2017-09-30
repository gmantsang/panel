package shards

import "github.com/flosch/pongo2"

// Templates contains the templates for shards
type Templates struct {
	List   *pongo2.Template
	Report *pongo2.Template
}

// CompileTemplates builds the pongo2 templates for shards
func CompileTemplates() (*Templates, error) {
	list, err := pongo2.FromFile("templates/shards/list.html")
	if err != nil {
		return nil, err
	}

	report, err := pongo2.FromFile("templates/shards/report.html")
	if err != nil {
		return nil, err
	}

	return &Templates{
		List:   list,
		Report: report,
	}, nil
}
