package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type downloadResult struct {
	Page  []byte
	Error error
}

type releaseView struct {
	DownloadURL string
	Version     string
	Md5         string
	Path        string
}

var rePackageURL = regexp.MustCompile(`(?i)<a href=\"(?P<url>.+?)#md5=.+?\">(?P<filename>.+?)</a>`)
var reDownloadURL = regexp.MustCompile(`(?i)<a href=\"(?P<url>.+?)\"\s+rel=\"download">(?P<version>.+?) download_url</a>`)

func simpleIndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering simple index")
	renderTemplate(w, "index", nil)
}

func getPage(url string) downloadResult {
	log.Println("Downloading page: " + url)
	res, err := http.Get(url)
	if err != nil {
		return downloadResult{Page: []byte(""), Error: err}
	}
	page, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return downloadResult{Page: []byte(""), Error: err}
	}
	return downloadResult{Page: page, Error: err}
}

func updateProxyCache(pkg Package) error {
	log.Println("Updating proxy cache for: " + pkg.Name)
	url := PypiMirror + "/simple/" + pkg.Name + "/"

	result := getPage(url)

	if result.Error != nil {
		// We'll try next time
		return result.Error
	}

	finalizeCache(pkg, result.Page)
	return nil
}

func finalizeCache(pkg Package, data []byte) {
	log.Println("Finalizing cache for: " + pkg.Name)
	returnData := string(data)
	// Replace the local package links with links to a local proxy
	// so we can cache that result as well.
	packageUris := rePackageURL.FindAllSubmatch(data, -1)
	for _, line := range packageUris {
		uri := line[1]
		filename := line[2]
		// TODO: This most certainly has edge cases that aren't addressed here.
		// Specifically: ".tar.gz" can't be the only filetype uploaded...
		versionSplit := strings.Split(string(uri), "-")
		almostVersion := versionSplit[len(versionSplit)-1]
		version := strings.Replace(almostVersion, ".tar.gz", "", -1)
		quoteduri := url.QueryEscape(PypiMirror + "/a/b/" + string(uri))
		replaceuri := "/fondu/cached-file/" + string(filename) + "?package=" + pkg.Name + "&release=" + version + "&original=" + quoteduri + "&name=" + url.QueryEscape(string(filename))
		returnData = strings.Replace(returnData, string(uri), replaceuri, -1)
	}

	// Replace the download links with links to a local proxy so that
	// we can cache the downloads as well.
	downloadUrls := reDownloadURL.FindAllSubmatch(data, -1)
	for _, line := range downloadUrls {
		uri := line[1]
		version := line[2]
		filename := pkg.Name + "-" + string(version) + ".tar.gz"
		quotedURI := url.QueryEscape(string(uri))
		replaceURI := "/fondu/cached-file/" + filename + "?package=" + pkg.Name + "&release=" + string(version) + "&original=" + quotedURI + "&name=" + url.QueryEscape(filename)
		returnData = strings.Replace(returnData, string(uri), replaceURI, -1)
	}

	pkg.SetProxy([]byte(returnData))
}

func renderProxy(w http.ResponseWriter, pkg Package) {
	log.Println("Rendering proxy for: " + pkg.Name)
	data := string(pkg.ProxyData())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Connection", "keep-alive")
	io.WriteString(w, data)
}

func buildReleaseMap(pkg Package) []releaseView {
	releaseMap := []releaseView{}
	for _, rel := range pkg.Releases() {
		releaseMap = append(releaseMap, releaseView{
			DownloadURL: rel.DownloadURL(),
			Version:     rel.Version,
			Md5:         rel.Md5(),
			Path:        rel.Path(),
		})
	}
	return releaseMap
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SIMPLE: " + r.URL.Path)
	paths := strings.Split(r.URL.Path, "/")

	if len(paths) == 3 && paths[len(paths)-1] == "" {
		simpleIndexHandler(w, r)
		return
	}

	if len(paths) == 3 && paths[len(paths)-1] != "" {
		http.NotFound(w, r)
		return
	}

	if len(paths) == 4 && paths[len(paths)-1] != "" {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		return
	}

	if len(paths) == 5 && paths[len(paths)-1] == "" {
		http.NotFound(w, r)
		return
	}

	name := paths[2]
	pkg := Package{Name: name, DataDir: FonduData}

	// The package is ours, so we serve it ourselves.
	if pkg.Exists() && !pkg.Proxied() {
		log.Print("Private package: " + name + ". Serving it.")
		releaseMap := buildReleaseMap(pkg)
		renderTemplate(w, "single", &releaseMap)
		return
	}

	// Public package, so just render the proxy
	if pkg.Proxied() {
		log.Print("Proxied package: " + name + ". Sending cached data.")
		go updateProxyCache(pkg)
		renderProxy(w, pkg)
		return
	}

	if err := updateProxyCache(pkg); err != nil {
		http.Error(w, "Unable to update proxy", http.StatusInternalServerError)
		return
	}

	renderProxy(w, pkg)
}
