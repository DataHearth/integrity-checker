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

var algorithmCommands = map[string][]string{
	"sha1":      {"shasum", "1"},
	"sha224":    {"shasum", "224"},
	"sha256":    {"shasum", "256"},
	"sha384":    {"shasum", "384"},
	"sha512":    {"shasum", "512"},
	"sha512224": {"shasum", "512224"},
	"sha512256": {"shasum", "512256"},
	"md5":       {"md5sum"},
}

var verbose bool

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

	if len(*checkFile) > 0 {
		CheckWithFile(*checkFile, *algorithm)
	} else if len(*filePath) > 0 && len(*checksum) > 0 {
		CheckWithChecksum(*checksum, *filePath, *algorithm)
	}
}

// CheckWithChecksum use a checksum chain and compare it to the result of the choosen algorithm.
// c is checksum chain, f is the file path, a is the algorithm
func CheckWithChecksum(c string, f string, a string) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command(algorithmCommands[a][0], "-a", algorithmCommands[a][1], f)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Run()

	if stderr.Len() > 0 {
		// Check if this is a checksum error
		matched, err := regexp.Match(`(shasum|md5sum): `+ f +`:`, stderr.Bytes())
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

	splittedFilePath := strings.Split(f, "/") // ! Linux File system. Need to add windows in the future
	splittedOutput := strings.Split(stdout.String(), " ")

	if splittedOutput[1] != splittedOutput[len(splittedFilePath)] {
		color.Red("Internal error")
		if verbose {
			fmt.Printf("The computed file did not match the file provided")
		}

		os.Exit(0)
	}

	if splittedOutput[0] != c {
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
func CheckWithFile(cf string, a string) {
	var stderr, stdout bytes.Buffer
	var cmd *exec.Cmd

	// execute the validation command with the right algorithm and check file
	if a == "md5" {
		// md5 algorithm
		cmd = exec.Command(algorithmCommands[a][0], "--check", cf)
	} else {
		// sha algorithm
		cmd = exec.Command(algorithmCommands[a][0], "-a", algorithmCommands[a][1], "--check", cf)
	}
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
