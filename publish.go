package main

import (
	"encoding/json"
	"log"
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
	log.Println("Registering: " + name + "-" + version)
	pkg := Package{Name: name, DataDir: FonduData}

	if pkg.Proxied() {
		pkg.Delete()
	}

	rel := Release{Name: name, Version: version, DataDir: pkg.DataDir}
	r.ParseForm()
	jsonData, _ := json.Marshal(r.Form)
	err := rel.StoreMetadata(jsonData)
	if err != nil {
		log.Println("Failed to write metadata to: " + rel.MetadataFile())
	}
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	name := r.FormValue("name")
	version := r.FormValue("version")
	log.Println("Uploading: " + name + "-" + version)
	content, _, _ := r.FormFile("content")

	if content == nil {
		http.Error(w, "No content provided", http.StatusBadRequest)
		return
	}

	pkg := Package{Name: name, DataDir: FonduData}
	if pkg.Proxied() {
		pkg.Delete()
	}

	rel := Release{Name: name, Version: version, DataDir: FonduData, Filename: name + "-" + version + ".tar.gz"}
	r.ParseForm()
	jsonData, _ := json.Marshal(r.Form)
	rel.StoreMetadata(jsonData)
	rel.StoreRelease(content)
}
