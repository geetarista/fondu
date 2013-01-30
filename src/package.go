package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Package struct {
	Name    string
	DataDir string
}

func (p Package) Releases() (releases []Release) {
	if _, err := os.Stat(p.Directory()); os.IsNotExist(err) {
		return make([]Release, 0)
	}

	versions, _ := filepath.Glob(p.Directory() + "/*")

	for _, v := range versions {
		stat, _ := os.Stat(v)
		filename := stat.Name()
		if string(filename[0]) == "." {
			continue
		}

		if stat.IsDir() == true {
			releases = append(releases, p.Release(filename))
		}
	}

	return releases
}

func (p Package) Directory() string {
	return filepath.Join(p.DataDir, p.Name)
}

func (p Package) Release(version string) Release {
	return Release{
		Name: p.Name,
		Version: version,
		DataDir: p.DataDir,
		Filename: p.Name + "-" + version + ".tar.gz",
	}
}

func (p Package) Exists() bool {
	return len(p.Releases()) > 0
}

func (p Package) ProxyFile() string {
	return filepath.Join(p.Directory(), "proxied")
}

func (p Package) Proxied() bool {
	_, err := os.Stat(p.ProxyFile())
	return err == nil
}

func (p Package) ProxyData() []byte {
	if !p.Proxied() {
		return make([]byte, 0)
	}

	data, _ := ioutil.ReadFile(p.ProxyFile())
	return data
}

func (p Package) SetProxy(data []byte) (err error) {
	dir := p.Directory()
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(p.ProxyFile(), data, 0644)
	return err
}

func (p Package) Delete() error {
	return os.RemoveAll(p.Directory())
}
