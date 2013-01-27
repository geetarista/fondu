package main

import (
	"encoding/json"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	switch r.FormValue(":action") {
	case "submit":
		registerHandler(w, r)
	case "file_upload":
		fileUploadHandler(w, r)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	version := r.FormValue("version")
	pkg := Package{Name: name}
	dataDir := Config.DataDir

	if pkg.Proxied() {
		pkg.Delete()
	}

	rel := Release{Name: name, Version: version, DataDir: dataDir}
	r.ParseForm()
	jsonData, _ := json.Marshal(r.Form)
	err := rel.StoreMetadata(jsonData)
	if err != nil {
		println("Failed to write metadata to: " + rel.MetadataFile())
	}
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	name := r.FormValue("name")
	version := r.FormValue("version")
	content, _, _ := r.FormFile("content")

	if content == nil {
		http.Error(w, "No content provided", http.StatusBadRequest)
	}

	pkg := Package{Name: name, DataDir: Config.DataDir}
	if pkg.Proxied() {
		pkg.Delete()
	}

	rel := Release{Name: name, Version: version, DataDir: Config.DataDir, Filename: name + "-" + version + ".tar.gz"}
	r.ParseForm()
	jsonData, _ := json.Marshal(r.Form)
	rel.StoreMetadata(jsonData)
	rel.StoreRelease(content)
}
