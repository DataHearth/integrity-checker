package utils

import "errors"

var (
	ErrCheckfile        = errors.New("integrity: please, provive a checkfile or a file and a checksum")
	ErrPathRequired     = errors.New("integrity: a file path is required")
	Algorithms = [8]string{
		"sha1",
		"sha224",
		"sha256",
		"sha384",
		"sha512",
		"sha512224",
		"sha512256",
		"md5",
	}
)