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

func main() {
	// Scope variables
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Setup user input
	verbose := flag.Bool("verbose", true, "Should the program be verbose \nDefault: true")
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
		var cmd *exec.Cmd

		// execute the validation command with the right algorithm and check file
		if *algorithm == "md5" {
			// md5 algorithm
			cmd = exec.Command(algorithmCommands[*algorithm][0], "--check", *checkFile)
		} else {
			// sha algorithm
			cmd = exec.Command(algorithmCommands[*algorithm][0], "-a", algorithmCommands[*algorithm][1], "--check", *checkFile)
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
				if *verbose {
					fmt.Print(stderr.String())
				}

				os.Exit(0)
			}
		}

		color.Green("Validation OK")
		if *verbose == true {
			fmt.Print(stdout.String())
		}

		os.Exit(0)
	} else if len(*filePath) > 0 && len(*checksum) > 0 {
		cmd := exec.Command(algorithmCommands[*algorithm][0], "-a", algorithmCommands[*algorithm][1], *filePath)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		cmd.Run()

		if stderr.Len() > 0 {
			// Check if this is a checksum error
			matched, err := regexp.Match(`(shasum|md5sum): `+*filePath+`:`, stderr.Bytes())
			if err != nil {
				log.Fatal(err)
			}

			if matched {
				color.Red("Validation error:")
				fmt.Print(stdout.String())

				os.Exit(0)
			} else {
				color.Red("Internal error")
				if *verbose {
					fmt.Print(stderr.String())
				}

				os.Exit(0)
			}
		}

		splittedFilePath := strings.Split(*filePath, "/") // ! Linux File system. Need to add windows in the future
		splittedOutput := strings.Split(stdout.String(), " ")

		if splittedOutput[1] != splittedOutput[len(splittedFilePath)] {
			color.Red("Internal error")
			if *verbose {
				fmt.Printf("The computed file did not match the file provided")
			}

			os.Exit(0)
		}

		if splittedOutput[0] != *checksum {
			color.Red("Validation error:")
			fmt.Printf("The checksum provided did not match the sha checksum: \n%s != %s\n", *checksum, splittedOutput[0])

			os.Exit(0)
		}

		color.Green("Validation OK")
		if *verbose == true {
			fmt.Print(stdout.String())
		}

		os.Exit(0)
	}
}
