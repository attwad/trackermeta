package html

import (
	"html/template"
	"io"

	"github.com/attwad/trackermeta/data"
)

// RenderHTML outputs HTML to display the tracker metadata in a useful way.
func RenderHTML(tf []data.TrackerFile, w io.Writer) error {
	tmpl := template.Must(template.ParseFiles("html/index.html"))
	return tmpl.Execute(w, tf)
}
