package main

import (
	"flag"
	"log"

	"github.com/datahearth/integrity-checker"
	"github.com/datahearth/integrity-checker/pkg/utils"
	"github.com/fatih/color"
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
		log.Fatalln(utils.ErrPathRequired.Error())
	}

	checker := checker.NewChecker(file, algorithm, checkfile, verbose)

	// checksum could be equal to "" if config.checkFile is true
	b, err := checker.IsValid(checksum)
	if err != nil {
		color.Red("%v", err)
	}

	if b {
		color.Green("Validation: OK")
	} else {
		color.Red("Validation: FAILED")
	}
}
