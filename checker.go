package main

import (
	"io"
	"os"
	"strings"

	"github.com/DataHearth/integrity-checker/utils"
)

// Checker is the base structure
type Checker struct {
	verbose bool
	config  CheckerConfig
}

// NewChecker check if required field are not missing
func NewChecker(c CheckerConfig, arg ...bool) Checker {
	if arg[0] {
		return Checker{verbose: true, config: c}
	}

	return Checker{verbose: false, config: c}
}

func (c Checker) isValid(ch string) (bool, error) {
	if !c.config.checkfile && len(ch) == 0 {
		return false, ErrCheckfile
	}

	f, err := os.Open(c.config.file)
	if err != nil {
		return false, err
	}

	defer f.Close()

	if len(ch) > 0 {
		_, err = io.Copy(c.config.algorithm, f)
		if err != nil {
			return false, err
		}

		return utils.Compare(ch, c.config.algorithm), nil
	}

	content := utils.SplitCheckfile(f)
	for _, line := range content {
		c.config.algorithm.Reset()

		lineContent := strings.Split(line, "  ")

		f, err := os.Open(lineContent[1])
		if err != nil {
			return false, err
		}

		defer f.Close()

		_, err = io.Copy(c.config.algorithm, f)
		if err != nil {
			return false, err
		}

		if !utils.Compare(lineContent[0], c.config.algorithm) {
			return false, nil
		}
	}

	return true, nil
}
