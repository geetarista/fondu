package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Package represents the high-level concept of a Python package.
type Package struct {
	Name    string
	DataDir string
}

// Releases finds all of the known releases for a package.
func (p Package) Releases() (releases []Release) {
	if _, err := os.Stat(p.Directory()); os.IsNotExist(err) {
		return make([]Release, 0)
	}

	versions, _ := filepath.Glob(p.Directory() + "/*")

	for _, v := range versions {
		stat, _ := os.Stat(v)
		filename := stat.Name()

		if stat.IsDir() == true {
			releases = append(releases, p.Release(filename))
		}
	}

	return releases
}

// Directory is the full path to where this package resides.
func (p Package) Directory() string {
	return filepath.Join(p.DataDir, p.Name)
}

// Release returns information about a specific version of a package.
func (p Package) Release(version string) Release {
	return Release{
		Name:     p.Name,
		Version:  version,
		DataDir:  p.DataDir,
		Filename: p.Name + "-" + version + ".tar.gz",
	}
}

// Exists returns whether a package exists.
func (p Package) Exists() bool {
	return len(p.Releases()) > 0
}

// ProxyFile is the path to the proxied file.
func (p Package) ProxyFile() string {
	return filepath.Join(p.Directory(), "proxied")
}

// Proxied returns whether a package is proxied.
func (p Package) Proxied() bool {
	_, err := os.Stat(p.ProxyFile())
	return err == nil
}

// ProxyData reads everything in the proxied file.
func (p Package) ProxyData() []byte {
	if !p.Proxied() {
		return make([]byte, 0)
	}

	data, _ := ioutil.ReadFile(p.ProxyFile())

	return data
}

// SetProxy will fill the proxied file with appropriate data.
func (p Package) SetProxy(data []byte) (err error) {
	dir := p.Directory()

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}

	if err != nil {
		return err
	}

	return ioutil.WriteFile(p.ProxyFile(), data, 0644)
}

// Delete removes everything related to a package.
func (p Package) Delete() error {
	return os.RemoveAll(p.Directory())
}
