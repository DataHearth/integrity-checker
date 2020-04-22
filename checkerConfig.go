package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"log"
)

// CheckerConfig defines the checker configuration
type CheckerConfig struct {
	file      string
	algorithm hash.Hash
	checkfile bool
}

// NewCheckerConfig returns a CheckerConfig object configured with function arguments.
// f is the file path, a is the algorithm and the fird is if the file is a checkfile
func NewCheckerConfig(f string, a string, arg ...bool) CheckerConfig {
	var algo hash.Hash

	if len(f) == 0 {
		log.Fatalln(ErrPathRequired.Error())
	}

	switch a {
	case Algorithms[0]:
		algo = sha1.New()
	case Algorithms[1]:
		algo = sha256.New224()
	case Algorithms[2]:
		algo = sha256.New()
	case Algorithms[3]:
		algo = sha512.New384()
	case Algorithms[4]:
		algo = sha512.New()
	case Algorithms[5]:
		algo = sha512.New512_224()
	case Algorithms[6]:
		algo = sha512.New512_256()
	case Algorithms[7]:
		algo = md5.New()
	default:
		algo = sha1.New()
	}

	if arg[0] {
		return CheckerConfig{algorithm: algo, checkfile: true, file: f}
	}

	return CheckerConfig{algorithm: algo, checkfile: false, file: f}
}