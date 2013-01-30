package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestRootGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(rootHandler))
	defer server.Close()

	res, err := http.Get(server.URL+"/")
	failIfError(t, err)

	assertEqual(t, "root get", res.StatusCode, http.StatusMethodNotAllowed)
}

func TestRegisterHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(rootHandler))
	defer server.Close()

	values := url.Values{
		":action": []string{"submit"},
		"name":    []string{"test"},
		"version": []string{"1.0.0"},
	}

	_, err := http.PostForm(server.URL+"/", values)
	failIfError(t, err)

	metadata, err := ioutil.ReadFile(Config.DataDir + "/test/1.0.0/metadata.json")

	assertEqual(t, "register handler", string(metadata), `{":action":["submit"],"name":["test"],"version":["1.0.0"]}`)
}

func TestRegisterHandlerProxied(t *testing.T) {
	proxyPackage.SetProxy(testBytes)
	server := httptest.NewServer(http.HandlerFunc(rootHandler))
	defer server.Close()

	values := url.Values{
		":action": []string{"submit"},
		"name":    []string{"proxypackage"},
		"version": []string{"1.0.0"},
	}

	_, err := http.PostForm(server.URL+"/", values)
	failIfError(t, err)

	metadata, err := ioutil.ReadFile(Config.DataDir + "/proxypackage/1.0.0/metadata.json")

	assertEqual(t, "register handler proxied", string(metadata), `{":action":["submit"],"name":["proxypackage"],"version":["1.0.0"]}`)
}

func TestFileUploadHandlerSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(rootHandler))
	defer server.Close()

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField(":action", "file_upload")
	w.WriteField("name", "proxypackage")
	w.WriteField("version", "1.0.0")
	w.WriteField("md5_digest", "abc123")
	fw, _ := w.CreateFormFile("content", "test.tar.gz")
	f, err := os.Open("./test/test.tar.gz")
	failIfError(t, err)
	defer f.Close()
	io.Copy(fw, f)
	ct := w.FormDataContentType()
	w.Close()

	res, err := http.Post(server.URL+"/submit", ct, buf)
	failIfError(t, err)
	assertEqual(t, "upload success", res.StatusCode, http.StatusOK)
}

func TestFileUploadHandlerProxied(t *testing.T) {
	proxyPackage.SetProxy(testBytes)
	server := httptest.NewServer(http.HandlerFunc(rootHandler))
	defer server.Close()

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField(":action", "file_upload")
	w.WriteField("name", "proxypackage")
	w.WriteField("version", "1.0.0")
	w.WriteField("md5_digest", "abc123")
	fw, _ := w.CreateFormFile("content", "test.tar.gz")
	f, err := os.Open("./test/test.tar.gz")
	failIfError(t, err)
	defer f.Close()
	io.Copy(fw, f)
	ct := w.FormDataContentType()
	w.Close()

	res, err := http.Post(server.URL+"/submit", ct, buf)
	failIfError(t, err)
	assertEqual(t, "upload proxied", res.StatusCode, http.StatusOK)
}

func TestFileUploadHandlerEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(rootHandler))
	defer server.Close()

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField(":action", "file_upload")
	w.WriteField("name", "proxypackage")
	w.WriteField("version", "1.0.0")
	w.WriteField("md5_digest", "abc123")
	ct := w.FormDataContentType()
	w.Close()

	res, err := http.Post(server.URL+"/submit", ct, buf)
	failIfError(t, err)
	assertEqual(t, "upload empty", res.StatusCode, http.StatusBadRequest)
}
