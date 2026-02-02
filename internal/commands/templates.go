package commands

import (
	"fmt"

	"github.com/promptforge/promptforge/internal/templates"
)

// ListTemplates prints available templates.
func ListTemplates() error {
	list := templates.List()
	if len(list) == 0 {
		fmt.Println("No templates available")
		return nil
	}

	for _, tmpl := range list {
		fmt.Printf("%s - %s\n", tmpl.Name, tmpl.Description)
	}

	return nil
}
