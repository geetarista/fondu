package main

import (
	"flag"
	"log"
	"mime"
	"net/http"
	"os"
)

// Fondu's version
const VERSION = "0.1.0"

var (
	// FonduData is the directory where fondu's data will be stored
	FonduData string
	// FonduPort is the port that fondu will listen on
	FonduPort string
	// PypiMirror is the mirror to pull new packages from
	PypiMirror string
)

func init() {
	flag.StringVar(&FonduData, "d", "data", "directory to save fondu data")
	flag.StringVar(&FonduPort, "p", "3638", "port for fondu to listen on")
	flag.StringVar(&PypiMirror, "m", "http://pypi.python.org", "pypi mirror")
	version := flag.Bool("v", false, "prints current fondu version")
	flag.Parse()

	if *version {
		os.Stdout.WriteString(VERSION + "\n")
		os.Exit(0)
	}
}

func main() {
	mime.AddExtensionType(".gz", "application/x-gzip")
	http.HandleFunc("/simple", simpleIndexHandler)
	http.HandleFunc("/simple/", simpleHandler)
	http.Handle("/file/", http.StripPrefix("/file", http.FileServer(http.Dir(FonduData))))
	http.HandleFunc("/fondu/cached-file/", cachedFileHandler)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("*", http.NotFound)
	log.Println("Starting fondu on port " + FonduPort)
	http.ListenAndServe(":"+FonduPort, nil)
}
