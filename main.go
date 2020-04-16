package main

import (
	"github.com/fatih/color"
	"regexp"
	"fmt"
	"flag"
	"os"
	"os/exec"
	"log"
	"bytes"
)

var algorithmCommands = map[string][]string {
	"sha1": {"shasum", "1"},
	"sha224": {"shasum", "224"},
	"sha256": {"shasum", "256"},
	"sha384": {"shasum", "384"},
	"sha512": {"shasum", "512"},
	"sha512224": {"shasum", "512224"},
	"sha512256": {"shasum", "512256"},
	"md5": {"md5sum"},
}

func main()  {
	// Setup user input
	verbose := flag.Bool("verbose", true, "Should the program be verbose \nDefault: true")
	filePath := flag.String("file", "", "Path the file to check integrity \nRequired")
	checkFile := flag.String("check", "", "Check file with file reference")
	algorithm := flag.String("algorithm", "sha256", "Algorithm to use (sha[1, 224, 256, 384, 512, 512224, 512224], md5 NOT RECOMMANDED). \nDefault: sha256")

	// Retrieve user input
	flag.Parse()

	// Verifyy if an argument is missing
	if len(*filePath) < 1 && len(*checkFile) < 1 {
		log.Println("The file path argument is required...");
		os.Exit(1)
	}

	if len(*checkFile) > 1 {
		var output bytes.Buffer
		var errorOutput bytes.Buffer
		var cmd *exec.Cmd
		
		// execute the validation command with the right algorithm and check file
		if *algorithm == "md5" {
			// md5 algorithm
			cmd = exec.Command(algorithmCommands[*algorithm][0], "--check", *checkFile)
		} else {
			// sha algorithm
			cmd = exec.Command(algorithmCommands[*algorithm][0], "-a", algorithmCommands[*algorithm][1], "--check", *checkFile)
		}
		cmd.Stdout = &output
		cmd.Stderr = &errorOutput

		// Run the command
		cmd.Run()

		if errorOutput.Len() > 0 {
			// Check if this is a checksum error
			matched, err := regexp.Match(`(shasum|md5sum): WARNING: [1-9]{1,9} computed checksum did NOT match`, errorOutput.Bytes())
			if err != nil {
				log.Fatal(err)
			}

			if matched {
				color.Red("Validation error:")
				fmt.Print(output.String())
	
				os.Exit(0)
			} else {
				color.Red("Internal error")
				if *verbose {
					fmt.Print(errorOutput.String())
				}
	
				os.Exit(0)
			}
		}
		
		color.Green("Validation OK")
		if *verbose == true {
			fmt.Print(output.String())
		}

		os.Exit(0)
	}
}