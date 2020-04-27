package main

import (
	"errors"
	"flag"
	"log"

	"github.com/fatih/color"
)

// Algorithms are all available algorithms and their associate commands
var Algorithms = [8]string{
	"sha1",
	"sha224",
	"sha256",
	"sha384",
	"sha512",
	"sha512224",
	"sha512256",
	"md5",
}

var (
	ErrCheckfile        = errors.New("integrity: please, provive a checkfile or a file and a checksum")
	ErrPathRequired     = errors.New("integrity: a file path is required")
)

func main() {
	var verbose bool
	var checkfile bool
	var file string
	var algorithm string
	var checksum string

	// Setup user input
	flag.BoolVar(&verbose, "verbose", true, "Should the program be verbose")
	flag.StringVar(&file, "file", "", "Path to your file (required)")
	flag.BoolVar(&checkfile, "check", true, "Is the file a checkfile")
	flag.StringVar(&algorithm, "algorithm", "sha1", "Supported algorithm (sha[1, 224, 256, 384, 512, 512224, 512224], md5 (NOT RECOMMANDED)).")
	flag.StringVar(&checksum, "checksum", "", "File's checksum (default empty)")

	// Retrieve user input
	flag.Parse()

	// Verify if an argument is missing
	if len(file) == 0 {
		log.Fatalln(ErrPathRequired.Error())
	}

	config := NewCheckerConfig(file, algorithm, checkfile)
	checker := NewChecker(config, verbose)

	// checksum could be equal to "" if config.checkFile is true
	b, err := checker.isValid(checksum)
	if err != nil {
		color.Red("%v", err)
	}

	if b {
		color.Green("Validation: OK")
	} else {
		color.Red("Validation: FAILED")
	}
}
