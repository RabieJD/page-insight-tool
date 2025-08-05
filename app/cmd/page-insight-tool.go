package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rabie/page-insight-tool/app/config"
	"github.com/rabie/page-insight-tool/app/router"
)

func main() {
	var configFile string
	var debug bool
	flag.StringVar(&configFile, "config", "", "config file location + name")
	flag.BoolVar(&debug, "debug", false, "debug log level (default Info)")
	flag.Parse()

	// Set up logging
	if debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// Load configuration
	cfg := config.LoadConfig(configFile)
	if cfg == nil {
		log.Fatal("Failed to load configuration")
	}

	// Create router
	r := router.New()

	// Start server
	log.Printf("Page Insight Tool listening on: %s", cfg.ServerAddress)
	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}
