package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type pageResult struct {
	Error error
}

func redirectToFile(w http.ResponseWriter, r *http.Request, release Release) {
	url := "/file/" + release.Path()
	log.Println("Redirecting to file: " + url)
	http.Redirect(w, r, url, http.StatusFound)
}

func downloadPage(url, file string) pageResult {
	log.Println("File doesn't exist yet. Downloading: " + url)
	res, err := http.Get(url)
	if err != nil {
		return pageResult{Error: err}
	}

	if file != "" {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return pageResult{Error: err}
		}
		err = ioutil.WriteFile(file, body, 0644)
		if err != nil {
			return pageResult{Error: err}
		}
	}
	return pageResult{Error: nil}
}

// If a file exists, just serve it directly. Otherwise, download, then serve.
// TODO: should we redirect to download_url if it exists?
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
		log.Println(release.Filename + "already exists. Redirecting to download.")
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
