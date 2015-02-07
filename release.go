package main

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
)

var (
	reMetaDownloadURL = regexp.MustCompile(`"download_url":\s?"([^"]+?)"`)
	reMd5             = regexp.MustCompile(`"md5_digest":\s?"([^"]+?)"`)
)

// Release represents a specific, versioned release of a package.
type Release struct {
	Name     string
	Version  string
	DataDir  string
	Filename string
}

// Metadata reads and returns everything in the metadata file.
func (r Release) Metadata() (data []byte, err error) {
	metaFile := r.MetadataFile()

	if _, err := os.Stat(metaFile); os.IsNotExist(err) {
		return data, err
	}

	return ioutil.ReadFile(metaFile)
}

// MetadataFile returns the path to the metadata file.
func (r Release) MetadataFile() string {
	return filepath.Join(r.Directory(), "metadata.json")
}

// Path represents the actual file for a release.
func (r Release) Path() string {
	return r.Name + "/" + r.Version + "/" + r.Filename
}

// DownloadURL fetches the location to download the package from.
func (r Release) DownloadURL() (url string) {
	mdata, err := r.Metadata()
	if err != nil {
		// panic(err.Error())
		return
	}

	matches := reMetaDownloadURL.FindSubmatch(mdata)

	if len(matches) > 1 {
		url = string(matches[1])
	}

	return
}

// Md5 returns the sum found from the original mirror.
func (r Release) Md5() (url string) {
	mdata, err := r.Metadata()
	if err != nil {
		// panic(err.Error())
		return
	}

	matches := reMd5.FindSubmatch(mdata)

	if len(matches) > 1 {
		url = string(matches[1])
	}

	return
}

// ReleaseFilePath points to the absolute location of a release.
func (r Release) ReleaseFilePath() string {
	return filepath.Join(r.Directory(), r.Filename)
}

// Exists returns whether a release is present.
func (r Release) Exists() bool {
	if r.Version == "" {
		return false
	}

	if _, err := os.Stat(r.ReleaseFilePath()); err == nil || r.DownloadURL() != "" {
		return true
	}

	return false
}

// StoreMetadata writes all metadata to a file.
func (r Release) StoreMetadata(data []byte) error {
	return ioutil.WriteFile(r.MetadataFile(), data, 0644)
}

// StoreRelease saves what was uploaded. Buffered in case of large files.
func (r Release) StoreRelease(data multipart.File) {
	f, err := os.Create(r.ReleaseFilePath())
	if err != nil {
		log.Println("Can't open file: " + r.ReleaseFilePath())
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)

	for {
		n, err := data.Read(buf)

		if err != nil && err != io.EOF {
			log.Println("Couldn't write file: " + r.ReleaseFilePath())
			break
		}

		if n == 0 {
			break
		}

		if _, err := f.Write(buf[:n]); err != nil {
			log.Println("Couldn't write file: " + r.ReleaseFilePath())
			break
		}
	}
}

// Directory is the path to the release directory,
// creating it if it doesn't exist.
func (r Release) Directory() string {
	projectDir := filepath.Join(r.DataDir, r.Name, r.Version)
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		os.MkdirAll(projectDir, 0755)
	}

	return projectDir
}
