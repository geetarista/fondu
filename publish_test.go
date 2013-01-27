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

func TestRegisterHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(registerHandler))
	defer server.Close()

	values := url.Values{
		"name":    []string{"test"},
		"version": []string{"1.0.0"},
	}

	_, err := http.PostForm(server.URL+"/submit", values)
	failIfError(t, err)

	metadata, err := ioutil.ReadFile(Config.DataDir + "/test/1.0.0/metadata.json")

	assertEqual(t, "register handler", string(metadata), `{"name":["test"],"version":["1.0.0"]}`)
}

func TestFileUploadHandlerSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(fileUploadHandler))
	defer server.Close()

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
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

	http.Post(server.URL+"/submit", ct, buf)
}
