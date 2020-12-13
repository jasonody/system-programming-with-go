// go run templates.go

package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	tpl, err := template.ParseGlob("./*.html")
	if err != nil {
		log.Fatal("Tempalte parse error:", err)
	}

	data := map[string]string{
		"name":       "Jin Kazama",
		"style":      "Karate",
		"appearance": "Tekken 3",
	}

	if err := tpl.Execute(os.Stdout, data); err != nil {
		log.Fatal("Template execute error:", err)
	}
}
