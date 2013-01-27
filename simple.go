package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type downloadResult struct {
	Page  []byte
	Error error
}

type releaseView struct {
	DownloadUrl string
	Version     string
	Md5         string
	Path        string
}

var rePackageUrl = regexp.MustCompile(`(?i)<a href=\"(?P<url>.+?)#md5=.+?\">(?P<filename>.+?)</a>`)
var reDownloadUrl = regexp.MustCompile(`(?i)<a href=\"(?P<url>.+?)\"\s+rel=\"download\">(?P<version>.+?) download_url</a>`)

func simpleIndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "simple/index", nil)
}

func getPage(url string) downloadResult {
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

func updateProxyCache(w http.ResponseWriter, pkg Package) error {
	url := Config.PypiMirror + "/simple/" + pkg.Name

	result := getPage(url)

	if result.Error != nil {
		// We'll try next time
		return result.Error
	}

	finalizeCache(w, pkg, result.Page)
	return nil
}

func finalizeCache(w http.ResponseWriter, pkg Package, data []byte) {
	returnData := string(data)
	// Replace the local package links with links to a local proxy
	// so we can cache that result as well.
	packageUris := rePackageUrl.FindAllSubmatch(data, -1)
	for _, line := range packageUris {
		uri := line[1]
		filename := line[2]
		// TODO: This most certainly has edge cases that aren't addressed here.
		// Specifically: ".tar.gz" can't be the only filetype uploaded...
		versionSplit := strings.Split(string(uri), "-")
		almostVersion := versionSplit[len(versionSplit)-1]
		version := strings.Replace(almostVersion, ".tar.gz", "", -1)
		quoteduri := url.QueryEscape(Config.PypiMirror + "/a/b/" + string(uri))
		replaceuri := "/fondu/cached-file/" + string(filename) + "?package=" + pkg.Name + "&release=" + version + "&original=" + quoteduri + "&name=" + url.QueryEscape(string(filename))
		println("Going to replace: " + string(uri))
		println("            with: " + replaceuri)
		returnData = strings.Replace(returnData, string(uri), replaceuri, -1)
	}

	// Replace the download links with links to a local proxy so that
	// we can cache the downloads as well.
	downloadUrls := reDownloadUrl.FindAll(data, -1)
	for _, line := range downloadUrls {
		uri := line[1]
		version := line[2]
		filename := pkg.Name + "-" + string(version) + ".tar.gz"
		quotedUri := url.QueryEscape(string(uri))
		replaceUri := "/fondu/cached-file/" + filename + "?package=" + pkg.Name + "&release=" + string(version) + "&original=" + quotedUri + "&name=" + url.QueryEscape(filename)
		returnData = strings.Replace(returnData, string(uri), replaceUri, -1)
	}

	pkg.SetProxy([]byte(returnData))
}

func renderProxy(w http.ResponseWriter, pkg Package) {
	io.WriteString(w, string(pkg.ProxyData()))
}

func buildReleaseMap(pkg Package) []releaseView {
	releaseMap := []releaseView{}
	for _, rel := range pkg.Releases() {
		metadata, err := rel.Metadata()
		if err != nil {
			break
		}
		md5Json, err := metadata.Get("md5_digest").Array()
		if err != nil {
			continue
		}
		md5 := md5Json[0]

		releaseMap = append(releaseMap, releaseView{
			DownloadUrl: rel.DownloadUrl(),
			Version:     rel.Version,
			Md5:         md5.(string),
			Path:        rel.Path(),
		})
	}
	return releaseMap
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	name := paths[len(paths)-2]
	pkg := Package{Name: name, DataDir: Config.DataDir}

	// The package is ours, so we serve it ourselves.
	if pkg.Exists() && !pkg.Proxied() {
		releaseMap := buildReleaseMap(pkg)
		renderTemplate(w, "simple/single", &releaseMap)
		return
	}

	// Public package, so just render the proxy
	if pkg.Proxied() {
		go updateProxyCache(w, pkg)
		renderProxy(w, pkg)
		return
	}

	if err := updateProxyCache(w, pkg); err != nil {
		http.Error(w, "Unable to update proxy", http.StatusInternalServerError)
	}

	renderProxy(w, pkg)
}
