package utils

import "errors"

var (
	// ErrCheckfile - occurs when a checkfile isn't provided
	ErrCheckfile = errors.New("integrity: please, provive a checkfile or a file and a checksum")
	// ErrPathRequired - occurs when filepath isn't provided
	ErrPathRequired = errors.New("integrity: a filepath is required")
	// ErrCopyAlgo - occurs when failed to copy file content into algorithm's hash for checking
	ErrCopyAlgo = errors.New("integrity: failed to copy data into algorithm")
	// ErrInvalid -
	ErrInvalid = errors.New("provided checksum isn't equal to file's checksum")
	// Algorithms - list of all supported algorithms
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
