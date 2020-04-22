package main

import (
	"errors"
	"flag"
	"log"
	"os"

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
	ErrUnknownAlgorithm = errors.New("integrity: please, provide a supported algorithm")
	ErrPathRequired = errors.New("integrity: a file path is required")
)

func main() {
	// Setup user input
	verbose := flag.Bool("verbose", true, "Should the program be verbose \nDefault: TRUE")
	file := flag.String("file", "", "Path to your file \nRequired")
	checkfile := flag.Bool("check", true, "Is the file a checkfile\nDefault: FALSE")
	algorithm := flag.String("algorithm", "sha1", "Supported algorithm (sha[1, 224, 256, 384, 512, 512224, 512224], md5 (NOT RECOMMANDED)). \nDefault: sha1")
	checksum := flag.String("checksum", "", "File's checksum \nDefault: EMPTY")

	// Retrieve user input
	flag.Parse()

	// Verify if an argument is missing
	if len(*file) == 0 {
		log.Fatalln(ErrPathRequired.Error())
	}

	config := NewCheckerConfig(*file, *algorithm, *checkfile)
	checker := NewChecker(config, *verbose)

	// checksum could be equal to "" if config.checkFile is true
	b, err := checker.isValid(*checksum)
	if err != nil {
		color.Red("%v", err)
	}

	if b {
		color.Green("Validation: OK")
	} else {
		color.Red("Validation: FAILED")
	}
}
