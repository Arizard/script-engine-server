package presenters

import (
	"bytes"
	"text/template"
	"fmt"
)

func renderTemplateToString(data interface{}, templatePath string) string {
	var buf bytes.Buffer
	t, err := template.ParseFiles(
		"presenters/templates/base.html",
		"presenters/templates/headerfooter.html",
		"presenters/templates/" + templatePath,
	)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	if err := t.Execute(&buf, data); err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	return buf.String()
}

// HTMLPresenter presents data in browser-renderable html.
type HTMLPresenter struct {}

func (HTMLPresenter) NotFound() string {
	return renderTemplateToString(nil, "404.html")
}


func (HTMLPresenter) InternalServerError() string {
	return renderTemplateToString(nil, "500.html")
}

func (HTMLPresenter) Index() string {
	data := struct {
		Something string
	}{
		Something: "Hello World",
	}
	return renderTemplateToString(data, "index.html")
}

// GetUserFile returns the HTML pretty-formatted user file view.
func (HTMLPresenter) GetUserFile(fileName string, publicURL string) string {
	data := struct {
		FileName string
		PublicURL string
	}{
		FileName: fileName,
		PublicURL: publicURL,
	}
	return renderTemplateToString(data, "getuserfile.html")
}