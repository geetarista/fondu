package main

import (
	"reflect"
	"regexp"
	"testing"
)

var validPackage = Package{Name: "validpackage", DataDir: Config.DataDir}
var validRelease = Release{
	Name:     validPackage.Name,
	Version:  "1.0.0",
	DataDir:  validPackage.DataDir,
	Filename: "test.tar.gz",
}
var metadata = `{
  "comment": "Test comment",
  "description": "Test description",
  "metadata_version": "1.0",
  "md5_digest": "bc06f913e06b680b6f021150a4e09044",
  "filetype": "sdist",
  "action": "file_upload",
  "pyversion": "",
  "keywords": [
      "test",
      "stuff",
      "things"
  ],
  "author_email": "test@example.com",
  "classifiers": [
      "Development Status :: 1 - Beta",
      "Intended Audience :: Developers",
      "License :: OSI Approved :: MIT License",
      "Operating System :: POSIX",
      "Programming Language :: Python",
      "Topic :: Software Development :: Libraries"
  ],
  "name": "test",
  "protcol_version": "1",
  "license": "MIT License",
  "author": "Test Person",
  "home_page": "https://github.com/test/test",
  "download_url": "http://example.com",
  "summary": "Just some test stuffs",
  "platform": "UNKNOWN",
  "version": "0.1.0"
}
`
var dummyPackage = Package{Name: "dummypackage", DataDir: Config.DataDir}
var dummyRelease = Release{Name: dummyPackage.Name, DataDir: dummyPackage.DataDir}
var proxyPackage = Package{Name: "proxypackage", DataDir: Config.DataDir}
var proxyRelease = Release{
	Name:     "proxypackage",
	Version:  "1.0.0",
	DataDir:  Config.DataDir,
	Filename: "proxypackage-1.0.0.tar.gz",
}
var privatePackage = Package{Name: "privatepackage", DataDir: Config.DataDir}
var privateRelease = Release{Name: privatePackage.Name, DataDir: privatePackage.DataDir}
var testBytes = []byte{'t', 'e', 's', 't', '\n'}

func failIfError(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}

func assertEqual(t *testing.T, name string, got, want interface{}) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%s => %v, want %v", name, got, want)
	}
}

func assertContains(t *testing.T, name, pattern, body string) {
	ok, err := regexp.MatchString(pattern, body)
	failIfError(t, err)
	if !ok {
		t.Errorf("%s: could not find %q in \n%q", name, pattern, body)
	}
}
