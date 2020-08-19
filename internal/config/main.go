package config

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"log"

	"github.com/datahearth/integrity-checker/internal/utils"
)

type Config struct {
	File      string
	Algorithm hash.Hash
	Checkfile bool
}

// NewConfig returns a CheckerConfig object configured with function arguments.
// f is the file path, a is the algorithm and the fird is if the file is a checkfile
func NewConfig(f string, a string, arg ...bool) Config {
	var algo hash.Hash

	if len(f) == 0 {
		log.Fatalln(utils.ErrPathRequired.Error())
	}

	switch a {
	case utils.Algorithms[0]:
		algo = sha1.New()
	case utils.Algorithms[1]:
		algo = sha256.New224()
	case utils.Algorithms[2]:
		algo = sha256.New()
	case utils.Algorithms[3]:
		algo = sha512.New384()
	case utils.Algorithms[4]:
		algo = sha512.New()
	case utils.Algorithms[5]:
		algo = sha512.New512_224()
	case utils.Algorithms[6]:
		algo = sha512.New512_256()
	case utils.Algorithms[7]:
		algo = md5.New()
	default:
		algo = sha1.New()
	}

	if arg[0] {
		return Config{Algorithm: algo, Checkfile: true, File: f}
	}

	return Config{Algorithm: algo, Checkfile: false, File: f}
}