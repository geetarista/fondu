package main

import (
	"testing"
)

var metadataTests = []struct {
	key string
	val interface{}
}{
	{"name", `"test"`},
	{"keywords", `["test","stuff","things"]`},
	{"comment", `"Test comment"`},
	{"download_url", `"http://example.com"`},
}

// func TestMetadata(t *testing.T) {
// 	m, err := validRelease.Metadata()
// 	failIfError(t, err)
//
// 	for _, mt := range metadataTests {
// 		got, err := m.Get(mt.key).MarshalJSON()
// 		failIfError(t, err)
//
// 		assertEqual(t, "metadata "+mt.key, mt.val, string(got))
// 	}
// }

func TestMetadataFile(t *testing.T) {
	got := validRelease.MetadataFile()
	want := validRelease.DataDir + "/" + validRelease.Name + "/" + validRelease.Version + "/metadata.json"

	assertEqual(t, "metadatafile", got, want)
}

func TestPath(t *testing.T) {
	got := validRelease.Path()
	want := validRelease.Name + "/" + validRelease.Version + "/" + validRelease.Filename

	assertEqual(t, "path", got, want)
}

func TestValidDownloadUrl(t *testing.T) {
	got := validRelease.DownloadUrl()
	want := "http://example.com"

	assertEqual(t, "valid download url", got, want)
}

func TestInvalidDownloadUrl(t *testing.T) {
	got := dummyRelease.DownloadUrl()
	want := ""

	assertEqual(t, "invalid download url", got, want)
}

func TestReleaseFilePath(t *testing.T) {
	got := validRelease.ReleaseFilePath()
	want := validRelease.Directory() + "/" + validRelease.Filename

	assertEqual(t, "release file path", got, want)
}

func TestValidExists(t *testing.T) {
	got := validRelease.Exists()

	assertEqual(t, "valid exists", got, true)
}

func TestInvalidExists(t *testing.T) {
	got := dummyRelease.Exists()

	assertEqual(t, "invalid exists", got, false)
}

func TestStoreMetadata(t *testing.T) {
	data := []byte(testMetadata)
	got := validRelease.StoreMetadata(data)

	assertEqual(t, "store metadata", got, nil)
}

// func TestStoreRelease(t *testing.T)

func TestReleaseDirectory(t *testing.T) {
	got := validRelease.Directory()
	want := validRelease.DataDir + "/" + validRelease.Name + "/" + validRelease.Version

	assertEqual(t, "project dir", got, want)
}
