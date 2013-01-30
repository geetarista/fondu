package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: Figure out if there's a way to assert the body. Right now we can't
// because httptest only sees one handler and it redirects to itself.
func TestSimpleIndexHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(simpleIndexHandler))
	defer server.Close()

	res, err := http.Get(server.URL + "/simple")
	failIfError(t, err)

	assertEqual(t, "simple index", res.StatusCode, 200)
}

func TestGetPage(t *testing.T) {
}

func TestSimpleSingleHandlerPrivate(t *testing.T) {
	privatePackage.SetProxy(testBytes)
	server := httptest.NewServer(http.HandlerFunc(simpleHandler))
	defer server.Close()

	res, err := http.Get(server.URL + "/simple/privatepackage/")
	failIfError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	failIfError(t, err)

	assertEqual(t, "simple single status", res.StatusCode, 200)
	assertContains(t, "simple single body", string(testBytes), string(body))
}

func TestSimpleSingleHandlerProxied(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(simpleHandler))
	defer server.Close()

	res, err := http.Get(server.URL + "/simple/proxypackage/")
	failIfError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	failIfError(t, err)

	assertEqual(t, "simple single status", res.StatusCode, 200)
	assertContains(t, "simple single body", `<a href="/file/proxypackage/1.0.0/proxypackage-1.0.0.tar.gz#md5=abc123">1.0.0</a>`, string(body))
}
