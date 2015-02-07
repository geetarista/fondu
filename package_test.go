package main

import (
	"os"
	"testing"
)

var testMetadata = `{
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

func TestSetupPackage(t *testing.T) {
	dummyPackage.Delete()
	err := validRelease.StoreMetadata([]byte(testMetadata))
	failIfError(t, err)
	validPackage.SetProxy(testBytes)
	f, err := os.Open("./test/test.tar.gz")
	failIfError(t, err)
	defer f.Close()
	validRelease.StoreRelease(f)
}

func TestEmptyReleases(t *testing.T) {
	got := dummyPackage.Releases()
	want := []Release{}
	assertEqual(t, "empty releases", got, want)
}

func TestValidReleases(t *testing.T) {
	got := validPackage.Releases()
	want := []Release{Release{Name: validPackage.Name, Version: "1.0.0", DataDir: validPackage.DataDir, Filename: "validpackage-1.0.0.tar.gz"}}
	assertEqual(t, "valid releases", got, want)
}

func TestProjectDirectory(t *testing.T) {
	got := validPackage.Directory()
	want := validPackage.DataDir + "/" + validPackage.Name

	assertEqual(t, "directory", got, want)
}

func TestRelease(t *testing.T) {
	want := Release{Name: validPackage.Name, Version: "1.0.0", DataDir: validPackage.DataDir, Filename: "validpackage-1.0.0.tar.gz"}
	got := validPackage.Release("1.0.0")

	assertEqual(t, "release", got, want)
}

func TestExists(t *testing.T) {
	want := true
	got := validPackage.Exists()

	assertEqual(t, "exists", got, want)
}

func TestProxyFile(t *testing.T) {
	got := validPackage.ProxyFile()
	want := validPackage.DataDir + "/" + validPackage.Name + "/proxied"

	assertEqual(t, "proxy file", got, want)
}

func TestProxied(t *testing.T) {
	got := validPackage.Proxied()
	want := true

	assertEqual(t, "proxied", got, want)
}

func TestProxyDataProxied(t *testing.T) {
	got := validPackage.ProxyData()
	want := testBytes

	assertEqual(t, "proxy data", got, want)
}

func TestProxyDataNotProxied(t *testing.T) {
	got := proxyPackage.ProxyData()
	want := []byte("")

	assertEqual(t, "proxy data", got, want)
}

func TestSetProxy(t *testing.T) {
	var err error
	proxyFile := validPackage.ProxyFile()
	err = os.Remove(proxyFile)
	assertEqual(t, "error", err, nil)

	err = validPackage.SetProxy(testBytes)
	assertEqual(t, "error", err, nil)
}

func TestDelete(t *testing.T) {
	var err error

	// Test that the directory does not exist and returns error
	err = dummyPackage.SetProxy(testBytes)
	assertEqual(t, "error", err, nil)

	// To create the dir
	dummyPackage.SetProxy(testBytes)

	err = dummyPackage.Delete()
	assertEqual(t, "error", err, nil)
}
