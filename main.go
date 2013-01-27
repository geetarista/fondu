package main

import (
	"flag"
	"github.com/kless/goconfig/config"
	"log"
	"net/http"
)

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
	flag.Parse()

	if configFile != "" {
		c, err := config.ReadDefault(configFile)
		if err != nil {
			log.Println("Unable to load config from: " + configFile)
			return
		}
		Config.DataDir, _ = c.String("fondu", "data_dir")
		Config.Port, _ = c.String("fondu", "port")
		Config.PypiMirror, _ = c.String("fondu", "pypi_mirror")
	}
}

func main() {
	http.HandleFunc("/simple", simpleIndexHandler)
	http.HandleFunc("/simple/", simpleHandler)
	http.Handle("/file/", http.StripPrefix("/file", http.FileServer(http.Dir(Config.DataDir))))
	http.HandleFunc("/fondu/cached-file/", cachedFileHandler)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("*", http.NotFound)
	log.Printf("Starting fondu on port " + Config.Port)
	http.ListenAndServe(":"+Config.Port, nil)
}
