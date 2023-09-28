package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Ometeor-Zheero-OMZ/go_web/pkg/config"
	"github.com/Ometeor-Zheero-OMZ/go_web/pkg/models"
)

var functions = template.FuncMap {}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a * config.AppConfig) {
	app = a
}

func AppDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
    var tc map[string]*template.Template
    if app.UseCache {
        // get the template cache from the app config
        tc = app.TemplateCache
    } else {
        tc, _ = CreateTemplateCache()
    }

    // get requested template from cache
    t, ok := tc[tmpl]
    if !ok {
        log.Fatal("Could not get template from template cache")
    }

    buf := new(bytes.Buffer)

    td = AppDefaultData(td)

    // Execute the template and write the result to the buffer
    err := t.Execute(buf, td)
    if err != nil {
        log.Println("Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Render the template
    _, err = buf.WriteTo(w)
    if err != nil {
        log.Println("Error writing template to browser:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}


func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page) // ts = template set
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts // it adds the final resulting template to my map
	}

	return myCache, nil
}