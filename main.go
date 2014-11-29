package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [run|configtest]\n", os.Args[0])
	fmt.Fprint(os.Stderr, "Configuration options:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func run(c Config) {
	configtest := (c.Args[0] == "configtest")

	checks, err := LoadChecks(c.CheckFilesPattern)
	if err != nil {
		log.Fatalf("Fatal error loading checks: %v\n", err)
	}

	if configtest {
		os.Exit(0)
	}

	log.Printf("Starting the server listening on %s.", c.BindAddress)
	RunServer(c, checks)
}

func main() {
	c := ParseConfig()

	if len(c.Args) == 1 {
		if c.Args[0] == "run" || c.Args[0] == "configtest" {
			run(c)
		} else {
			usage()
		}
	} else {
		usage()
	}
}
