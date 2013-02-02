package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type pageResult struct {
	Error error
}

func redirectToFile(w http.ResponseWriter, r *http.Request, release Release) {
	url := ""
	if release.DownloadUrl() != "" {
		println("DownloadUrl exists")
		url = release.DownloadUrl()
	} else {
		url = "/file/" + release.Path()
	}
	log.Println("Redirecting to: " + url)
	http.Redirect(w, r, url, http.StatusFound)
	return
}

func downloadPage(url, file string) pageResult {
	log.Println("File doesn't exist yet. Downloading: " + url)
	tr := &http.Transport{DisableCompression: true}
	client := &http.Client{Transport: tr}
	res, err := client.Get(url)

	if err != nil {
		log.Println("Error downloading: " + url)
		return pageResult{Error: err}
	}

	if file != "" {
		f, err := os.Create(file)
		if err != nil {
			log.Printf("Error creating file: " + file)
			return pageResult{Error: err}
		}
		_, err = io.Copy(f, res.Body)
		if err != nil {
			log.Printf("Error writing to: " + file)
			return pageResult{Error: err}
		}
	}
	return pageResult{Error: nil}
}

// If a file exists, just serve it directly. Otherwise, download, then serve.
func cachedFileHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	name := paths[len(paths)-1]
	log.Println("Request for cached file: " + name)
	release := Release{
		Name:     r.FormValue("package"),
		Version:  r.FormValue("release"),
		DataDir:  Config.DataDir,
		Filename: name,
	}

	if release.Exists() {
		log.Println(release.Filename + " already exists. Redirecting to download.")
		redirectToFile(w, r, release)
		return
	}

	original := r.FormValue("original")
	res := downloadPage(original, release.ReleaseFilePath())

	if res.Error != nil {
		http.Error(w, "Unable to download page", http.StatusInternalServerError)
		return
	}

	redirectToFile(w, r, release)
}
