package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/kaleanup-indx/hotel_reservation/pkg/config"
	"github.com/kaleanup-indx/hotel_reservation/pkg/models"
)

const tmplDirectoryName = "./templates/"
const tmplFileSuffix = ".page.tmpl"
const layoutFileSuffix = ".layout.tmpl"

var app *config.AppConfig

// NewTemplates set the config for the package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		//create template cache
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get requested template from cache
	t, ok := tc[tmpl+tmplFileSuffix]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	//create a buffer
	buf := new(bytes.Buffer)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	//render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)

	//same as line no 24. use this over that.
	myCache := map[string]*template.Template{}

	// get all of the files named *.pages.tmpl from templates
	pages, err := filepath.Glob(tmplDirectoryName + "*" + tmplFileSuffix)
	if err != nil {
		return myCache, err
	}

	//get all layout files
	matches, err := filepath.Glob(tmplDirectoryName + "*" + layoutFileSuffix)
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *.pages.tmpl
	for _, page := range pages {
		//get name of the file from path
		name := filepath.Base(page)

		//create template set and create a new template with the name of the file
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(tmplDirectoryName + "*" + layoutFileSuffix)
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
