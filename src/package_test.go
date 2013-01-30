package main

import (
	"os"
	"testing"
)

func TestSetupPackage(t *testing.T) {
	dummyPackage.Delete()
	validRelease.StoreMetadata([]byte(metadata))
	validPackage.SetProxy(testBytes)
	f, err := os.Open("./test/test.tar.gz")
	failIfError(t, err)
	defer f.Close()
	validRelease.StoreRelease(f)
}

func TestEmptyReleases(t *testing.T) {
	got := dummyPackage.Releases()
	want := make([]Release, 0)
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
