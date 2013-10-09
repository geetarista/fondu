package main

import (
	"log"
	"net/http"
	"text/template"
)

var templates = map[string]string{
	"base": `<html>
<head>
  <title>Fondu</title>
</head>
<body>
  {{ template "content" . }}
</body>
</html>`,
	"index": `{{ define "content" }}
<!-- TODO: List all packages -->
{{ end }}`,
	"single": `{{ define "content" }}
<ul>
{{ range . }}
{{ if .DownloadUrl }}
  <li><a href="{{ .DownloadUrl }}" rel="download">{{ .Version }} download</a></li>
{{ end }}
{{ if .Path }}
  <li><a href="/file/{{ .Path }}#md5={{ .Md5 }}">{{ .Version }}</a></li>
{{ end }}
{{ end }}
</ul>
{{ end }}`,
}

func renderTemplate(w http.ResponseWriter, tmpl string, tmplData interface{}) {
	log.Println("Rendering template: " + tmpl)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.New(tmpl).Parse(templates["base"] + templates[tmpl])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := t.Execute(w, tmplData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
