package main

import (
	"errors"
	"net/http"
	"os"
	"text/template"
	"time"
)

var tmplCache = make(map[string]*tmplInfo)

type tmplInfo struct {
	tmpl  *template.Template
	mtime int64
}

func checkTemplateExists(path string) (mtime int64, err error) {
	dir, err := os.Stat(path)
	if err != nil {
		return time.Now().Unix(), err
	}
	if dir.IsDir() {
		return time.Now().Unix(), errors.New("'" + path + "' is not a regular file")
	}

	return dir.ModTime().Unix(), nil
}

func loadTemplate(layout, tmpl string) (t *template.Template, err error) {
	key := layout + "-" + tmpl
	lmt, err := checkTemplateExists(layout)
	if err != nil {
		return nil, err
	}
	vmt, err := checkTemplateExists(tmpl)
	if err != nil {
		return nil, err
	}
	ti, _ := tmplCache[key]
	if ti == nil || lmt > ti.mtime || vmt > ti.mtime {

		t, err := template.ParseFiles(layout, tmpl)
		if err != nil {
			return nil, err
		}

		ti = &tmplInfo{
			mtime: time.Now().Unix(),
			tmpl:  t,
		}
		tmplCache[key] = ti
	}
	return ti.tmpl, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, tmplData interface{}) {
	l := "views/simple/base.html"
	v := "views/" + tmpl + ".html"
	t, err := loadTemplate(l, v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, tmplData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
