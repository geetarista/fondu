package main

import (
	simplejson "github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"mime/multipart"
)

type Release struct {
	Name     string
	Version  string
	DataDir  string
	Filename string
}

func (r Release) Metadata() (j *simplejson.Json, err error) {
	metaFile := r.MetadataFile()

	if _, err := os.Stat(metaFile); os.IsNotExist(err) {
		return j, err
	}

	data, err := ioutil.ReadFile(metaFile)

	if err != nil {
		return j, err
	}

	j, err = simplejson.NewJson(data)

	return j, err
}

func (r Release) MetadataFile() string {
	return filepath.Join(r.Directory(), "metadata.json")
}

func (r Release) Path() string {
	return r.Name + "/" + r.Version + "/" + r.Filename
}

func (r Release) DownloadUrl() (url string) {
	mdata, err := r.Metadata()
	if err != nil {
		return
	}

	if json, ok := mdata.CheckGet("download_url"); ok {
		url, _ = json.String()
		return
	}

	return
}

func (r Release) ReleaseFilePath() string {
	return filepath.Join(r.Directory(), r.Filename)
}

func (r Release) Exists() bool {
	if r.Version == "" {
		return false
	}

	if _, err := os.Stat(r.ReleaseFilePath()); err == nil || r.DownloadUrl() != "" {
		return true
	}

	return false
}

func (r Release) StoreMetadata(data []byte) error {
	return ioutil.WriteFile(r.MetadataFile(), data, 0644)
}

// Store the release that was uploaded. Buffered in case of large files.
func (r Release) StoreRelease(data multipart.File) {
	f, err := os.Create(r.ReleaseFilePath())
	if err != nil {
		log.Printf("Can't open file: " + r.ReleaseFilePath())
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)

	for {
		n, err := data.Read(buf)

		if err != nil && err != io.EOF {
			log.Printf("Couldn't write file: " + r.ReleaseFilePath())
			break
		}

		if n == 0 {
			break
		}

		if _, err := f.Write(buf[:n]); err != nil {
			log.Printf("Couldn't write file: " + r.ReleaseFilePath())
			break
		}
	}
}

func (r Release) Directory() string {
	projectDir := filepath.Join(r.DataDir, r.Name, r.Version)
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		os.MkdirAll(projectDir, 0755)
	}
	return projectDir
}
