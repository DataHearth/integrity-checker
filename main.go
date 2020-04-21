package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

// ErrCheckfile zefzef
var (
	ErrCheckfile        = errors.New("integrity: please, provive a checkfile or a file and a checksum")
	ErrUnknownAlgorithm = errors.New("integrity: please, provide a supported algorithm")
)

// CheckerConfig defines the checker configuration
type CheckerConfig struct {
	file      string
	algorithm string
	checkfile bool
}

// Checker fezfz
type Checker struct {
	verbose bool
	config  CheckerConfig
}

func main() {
	// Setup user input
	verbose := flag.Bool("verbose", true, "Should the program be verbose \nDefault: true")
	file := flag.String("file", "", "File to check (could also be a list of file check with related checksum) \nRequired")
	checkfile := flag.Bool("check", false, "Check file with file reference")
	algorithm := flag.String("algorithm", "sha256", "Algorithm to use (sha[1, 224, 256, 384, 512, 512224, 512224], md5 NOT RECOMMANDED). \nDefault: sha256")
	checksum := flag.String("checksum", "", "Checksum for the file")

	// Retrieve user input
	flag.Parse()

	// Verify if an argument is missing
	if len(*file) == 0 {
		log.Println("The file path argument is required...")

		os.Exit(1)
	}

	config := CheckerConfig{*file, *algorithm, *checkfile}
	checker := Checker{*verbose, config}

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

func (c Checker) isValid(ch string) (bool, error) {
	var h hash.Hash

	if !c.config.checkfile && len(ch) == 0 {
		return false, ErrCheckfile
	}

	f, err := os.Open(c.config.file)
	if err != nil {
		return false, os.ErrNotExist
	}

	defer f.Close()

	switch c.config.algorithm {
	case Algorithms[0]:
		h = sha1.New()
	case Algorithms[1]:
		h = sha256.New224()
	case Algorithms[2]:
		h = sha256.New()
	case Algorithms[3]:
		h = sha512.New384()
	case Algorithms[4]:
		h = sha512.New()
	case Algorithms[5]:
		h = sha512.New512_224()
	case Algorithms[6]:
		h = sha512.New512_256()
	case Algorithms[7]:
		h = md5.New()
	default:
		return false, ErrUnknownAlgorithm
	}

	if len(ch) == 0 {
		content := splitCheckfile(f)
		for _, line := range content {
			h.Reset()

			lineContent := strings.Split(line, "  ")

			f, err := os.Open(lineContent[1])
			if err != nil {
				return false, os.ErrNotExist
			}

			defer f.Close()

			_, err = io.Copy(h, f)
			if err != nil {
				return false, err
			}

			if !compare(lineContent[0], h) {
				return false, nil
			}
		}

		return true, nil
	}

	_, err = io.Copy(h, f)
	if err != nil {
		return false, err
	}

	return compare(ch, h), nil
}

func splitCheckfile(f *os.File) []string {
	d, err := ioutil.ReadAll(f)
	if err != nil {
		color.RedString(err.Error())
	}

	s := strings.Split(string(d), "\n")
	trim := s[:len(s)-1]

	return trim
}

func compare(c string, h hash.Hash) bool {
	d, err := hex.DecodeString(c)
	if err != nil {
		fmt.Println(err)
	}

	return bytes.Equal(d, h.Sum(nil))
}
