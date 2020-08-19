package checker

import (
	"io"
	"os"
	"strings"

	"github.com/datahearth/integrity-checker/internal/config"
	"github.com/datahearth/integrity-checker/internal/utils"
)

type checker struct {
	verbose bool
	config  config.Config
}

// NewChecker check if required field are not missing
func NewChecker(file, algorithm string, checkfile bool, arg ...bool) checker {
	config := config.NewConfig(file, algorithm, checkfile)
	if arg[0] {
		return checker{verbose: true, config: config}
	}

	return checker{verbose: false, config: config}
}

func (c checker) IsValid(ch string) (bool, error) {
	if !c.config.Checkfile && len(ch) == 0 {
		return false, utils.ErrCheckfile
	}

	f, err := os.Open(c.config.File)
	if err != nil {
		return false, err
	}

	defer f.Close()

	if len(ch) > 0 {
		_, err = io.Copy(c.config.Algorithm, f)
		if err != nil {
			return false, err
		}

		return utils.Compare(ch, c.config.Algorithm), nil
	}

	content := utils.SplitCheckfile(f)
	for _, line := range content {
		c.config.Algorithm.Reset()

		lineContent := strings.Split(line, "  ")

		f, err := os.Open(lineContent[1])
		if err != nil {
			return false, err
		}

		defer f.Close()

		_, err = io.Copy(c.config.Algorithm, f)
		if err != nil {
			return false, err
		}

		if !utils.Compare(lineContent[0], c.config.Algorithm) {
			return false, nil
		}
	}

	return true, nil
}
