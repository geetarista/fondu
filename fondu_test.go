package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSetupFondu(t *testing.T) {
	f, err := os.Open("./test/test.tar.gz")
	failIfError(t, err)
	defer f.Close()
	proxyRelease.StoreRelease(f)
}

func TestCachedFileHandlerExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(cachedFileHandler))
	defer server.Close()

	// Making a round-trip request to not follow the redirect
	tr := &http.Transport{}
	req, err := http.NewRequest("GET", server.URL+"/fondu/cached-file/proxypackage-1.0.0.tar.gz?package=proxypackage&release=1.0.0&original=test.com", nil)
	res, err := tr.RoundTrip(req)
	failIfError(t, err)

	assertEqual(t, "cached file handler exists", res.StatusCode, 302)
}

func TestCachedFileHandlerDownload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(cachedFileHandler))
	defer server.Close()

	// Making a round-trip request to not follow the redirect
	tr := &http.Transport{}
	c := &http.Client{Transport: tr}
	res, err := c.Get(server.URL + "/fondu/cached-file/test-1.0.0.tar.gz?package=test&release=1.0.0&original=/file/test/test.tar.gz")
	failIfError(t, err)

	assertEqual(t, "cached file handler download", res.StatusCode, http.StatusInternalServerError)
}
