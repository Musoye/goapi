package main

import (
	"log"
	"os"
	"text/template"
)

type Note struct {
	Title       string
	Description string
}

const tmpl = `Notes are:
{{range .}}
    Title: {{.Title}}
    DESCRIPTION: {{.Description}}
{{end}}`

func main() {
	notes := []Note{
		{"Title 1", "Description 1"},
		{"Title 2", "Description 2"},
		{"Title 3", "Description 3"},
	}

	t := template.Must(template.New("note").Parse(tmpl))
	if err := t.Execute(os.Stdout, notes); err != nil {
		log.Fatalf("t.Execute: %v", err)
	}
}
