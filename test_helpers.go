package main

import (
	"reflect"
	"regexp"
	"testing"
)

var validPackage = Package{Name: "validpackage", DataDir: "data"}
var validRelease = Release{
	Name:     validPackage.Name,
	Version:  "1.0.0",
	DataDir:  validPackage.DataDir,
	Filename: "test.tar.gz",
}
var dummyPackage = Package{Name: "dummypackage", DataDir: "data"}
var dummyRelease = Release{Name: dummyPackage.Name, DataDir: dummyPackage.DataDir}
var proxyPackage = Package{Name: "proxypackage", DataDir: "data"}
var proxyRelease = Release{
	Name:     "proxypackage",
	Version:  "1.0.0",
	DataDir:  "data",
	Filename: "proxypackage-1.0.0.tar.gz",
}
var privatePackage = Package{Name: "privatepackage", DataDir: "data"}
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
