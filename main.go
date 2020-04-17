package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// AlgorithmCommands are all available algorithms and their associate commands
var AlgorithmCommands = map[string][]string{
	"sha1":      {"shasum", "-a", "1"},
	"sha224":    {"shasum", "-a", "224"},
	"sha256":    {"shasum", "-a", "256"},
	"sha384":    {"shasum", "-a", "384"},
	"sha512":    {"shasum", "-a", "512"},
	"sha512224": {"shasum", "-a", "512224"},
	"sha512256": {"shasum", "-a", "512256"},
	"md5":       {"md5sum"},
}

var verbose bool

// Checker type with all informations
type Checker struct {
	filePath  string
	checkFile string
	algorithm string
	checksum  string
}

func main() {
	// Setup user input
	flag.BoolVar(&verbose, "verbose", true, "Should the program be verbose \nDefault: true")
	filePath := flag.String("file", "", "Path the file to check integrity \nRequired")
	checkFile := flag.String("check", "", "Check file with file reference")
	algorithm := flag.String("algorithm", "sha256", "Algorithm to use (sha[1, 224, 256, 384, 512, 512224, 512224], md5 NOT RECOMMANDED). \nDefault: sha256")
	checksum := flag.String("checksum", "", "Checksum for the file")

	// Retrieve user input
	flag.Parse()

	// Verifyy if an argument is missing
	if len(*filePath) == 0 && len(*checkFile) == 0 {
		log.Println("The file path argument is required...")

		os.Exit(1)
	}

	checker := Checker{*filePath, *checkFile, *algorithm, *checksum}
	if len(*checkFile) > 0 {
		checker.CheckWithFile()
	} else if len(*filePath) > 0 && len(*checksum) > 0 {
		checker.CheckWithChecksum()
	}
}

// CheckWithChecksum use a checksum chain and compare it to the result of the choosen algorithm.
// c is checksum chain, f is the file path, a is the algorithm
func (c Checker) CheckWithChecksum() {
	var stderr, stdout bytes.Buffer
	algorithmArgs := append(AlgorithmCommands[c.algorithm][1:], c.filePath)
	
	cmd := exec.Command(AlgorithmCommands[c.algorithm][0], algorithmArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Run()

	if stderr.Len() > 0 {
		// Check if this is a checksum error
		matched, err := regexp.Match(`(shasum|md5sum): `+ c.filePath +`:`, stderr.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		if matched {
			color.Red("Validation error:")
			fmt.Print(stdout.String())

			os.Exit(0)
		} else {
			color.Red("Internal error")
			if verbose {
				fmt.Print(stderr.String())
			}

			os.Exit(0)
		}
	}

	splittedFilePath := strings.Split(c.filePath, "/") // ! Linux File system. Need to add windows in the future
	splittedOutput := strings.Split(stdout.String(), " ")

	if splittedOutput[1] != splittedOutput[len(splittedFilePath)] {
		color.Red("Internal error")
		if verbose {
			fmt.Printf("The computed file did not match the file provided")
		}

		os.Exit(0)
	}

	if splittedOutput[0] != c.checksum {
		color.Red("Validation error:")
		fmt.Printf("The checksum provided did not match the sha checksum: \n%s != %s\n", c, splittedOutput[0])

		os.Exit(0)
	}

	color.Green("Validation OK")
	if verbose {
		fmt.Print(stdout.String())
	}

	os.Exit(0)
}

// CheckWithFile use a checksum file and use the algorithm provided
// cf is checksum file, a is the algorithm
func (c Checker) CheckWithFile() {
	var stderr, stdout bytes.Buffer

	algorithmArgs := append(AlgorithmCommands[c.algorithm][1:], "--check", c.checkFile)

	cmd := exec.Command(AlgorithmCommands[c.algorithm][0], algorithmArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	cmd.Run()

	if stderr.Len() > 0 {
		// Check if this is a checksum error
		matched, err := regexp.Match(`(shasum|md5sum): WARNING: [1-9]{1,9} computed checksum did NOT match`, stderr.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		if matched {
			color.Red("Validation error:")
			fmt.Print(stdout.String())

			os.Exit(0)
		} else {
			color.Red("Internal error")
			if verbose {
				fmt.Print(stderr.String())
			}

			os.Exit(0)
		}
	}

	color.Green("Validation OK")
	if verbose == true {
		fmt.Print(stdout.String())
	}

	os.Exit(0)
}
