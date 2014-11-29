package main

import (
	"flag"
	"fmt"
)

type Config struct {
	ElasticSearchURL      string
	ElasticSearchCountURL string
	BindAddress           string
	CheckFilesPattern     string
	RequestTimeout        int
	Args                  []string
}

func ParseConfig() (c Config) {
	c = Config{}

	// @TODO Validate these two URLs.
	flag.StringVar(&c.ElasticSearchURL, "esurl", "http://localhost:9200", "The URL to the ElasticSearch server that will be checked.")
	flag.StringVar(&c.BindAddress, "bindaddr", "0.0.0.0:8080", "The address that basecoach should bind to.")
	flag.StringVar(&c.CheckFilesPattern, "checkfiles", "/opt/basecoach/config/*.yml", "The path to search for check definition files. Globs are allowed.")
	flag.IntVar(&c.RequestTimeout, "timeout", 2000, "The timeout value for HTTP requests to ElasticSearch in milliseconds.")

	// Add the _count URL here, as in v1, all requests ElasticSearch _count queries.
	c.ElasticSearchCountURL = fmt.Sprintf("%s/_count", c.ElasticSearchURL)

	flag.Parse()
	c.Args = flag.Args()

	return c
}
