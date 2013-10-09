package main

import (
	"flag"
	"fmt"
	"github.com/robfig/config"
	"log"
	"mime"
	"net/http"
	"os"
)

const VERSION = "0.1.0"

var configFile string

type fonduConfig struct {
	DataDir    string
	Port       string
	PypiMirror string
}

var Config = fonduConfig{
	"data",
	"3638",
	"http://pypi.python.org",
}

// Parse flags and set up configuration
func init() {
	flag.StringVar(&configFile, "f", "", "path to a configuration file")
	version := flag.Bool("v", false, "prints current fondu version")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if configFile != "" {
		c, err := config.ReadDefault(configFile)
		if err != nil {
			log.Println("Unable to load config from: " + configFile)
			return
		}
		d, _ := c.String("fondu", "data_dir")
		if d != "" {
			Config.DataDir = d
		}
		p, _ := c.String("fondu", "port")
		if p != "" {
			Config.Port = p
		}
		m, _ := c.String("fondu", "pypi_mirror")
		if m != "" {
			Config.PypiMirror = m
		}
	}
}

func main() {
	mime.AddExtensionType(".gz", "application/x-gzip")
	http.HandleFunc("/simple", simpleIndexHandler)
	http.HandleFunc("/simple/", simpleHandler)
	http.Handle("/file/", http.StripPrefix("/file", http.FileServer(http.Dir(Config.DataDir))))
	http.HandleFunc("/fondu/cached-file/", cachedFileHandler)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("*", http.NotFound)
	log.Println("Starting fondu on port " + Config.Port)
	http.ListenAndServe(":"+Config.Port, nil)
}
