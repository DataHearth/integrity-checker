package checker

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/datahearth/integrity-checker/utils"
	"github.com/sirupsen/logrus"
)

// Check - default interface
type Check interface {
	CheckWithChecksum(checksum, filepath string) (bool, error)
	CheckWithCheckfile(filepath string) (bool, error)
}

// Checker - default structure
type Checker struct {
	verbose   bool
	algorithm hash.Hash
	logger    logrus.FieldLogger
	wg        sync.WaitGroup
}

// NewChecker - check if required field are not missing
func NewChecker(algorithm string, logger logrus.FieldLogger, verbose bool) Check {
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	var algo hash.Hash

	switch algorithm {
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

	return Checker{
		verbose:   verbose,
		algorithm: algo,
		logger:    logger,
		wg:        new(sync.WaitGroup),
	}
}

// CheckWithChecksum - check a single file with a checksum
func (c Checker) CheckWithChecksum(checksum, filepath string) (bool, error) {
	log := c.logger.WithFields(logrus.Fields{
		"method":   "CheckWithChecksum",
		"file":     filepath,
		"checksum": checksum,
	})

	log.Debug("reading file")
	// create a file reader
	file, err := os.Open(filepath)
	if err != nil {
		log.WithError(err).Error("failed to read file")
		return false, err
	}
	defer file.Close()

	log.Debug("copying file's data into algorithm to create a hash")
	// create file hash
	if _, err := io.Copy(c.algorithm, file); err != nil {
		log.WithError(err).Error("failed to create file hash")
		return false, utils.ErrCopyAlgo
	}

	log.Debug("comparing file hash with checksum")
	// check whether the file's hash is equal to the checksum
	if err := utils.Compare(checksum, c.algorithm); err != nil {
		log.Error(err.Error())
		return false, err
	}
	c.algorithm.Reset()

	return true, nil
}

// CheckWithCheckfile - check a list of files or directories with a custom checkfile
func (c Checker) CheckWithCheckfile(checkfilePath string) (bool, error) {
	log := c.logger.WithFields(logrus.Fields{
		"method": "CheckWithCheckfile",
		"file":   checkfilePath,
	})

	log.Debug("reading checkfile")
	// read checkfile data
	data, err := ioutil.ReadFile(checkfilePath)
	if err != nil {
		return false, err
	}

	log.Debug("splitting checkfile data into checklist")
	// split data into a map[checksum]file
	checklist := utils.SplitCheckfile(data)

	for checksum, path := range checklist {
		info, err := os.Stat(filepath.Join(checkfilePath, checksum))
		if !os.IsNotExist(err) {
			log.Info("file/folder path detected... Determine if we should iterate (folder only)")
			if !info.IsDir() {
				log.Info("path is not a folder. Skipping")
				continue
			}

			log.Debug("Adding a checkfile to the queue")
			go c.CheckWithCheckfile(filepath.Join(), checksum, "integrity_checkfile")
		}

		c.wg.Add(1)
		go func() {
			c.CheckWithChecksum(checksum, filepath.Join(checkfilePath, path))
			c.wg.Done()
		}()
	}
	// Clear this instance of waitgroup
	c.wg.Done()

	// Wait for all checkfile waitgroup to finish
	c.checkfileWG.Wait()
	return true, nil
}

// IsValid - check whether the provided checksum is valid
// func (c Checker) IsValid(checksum string) (bool, error) {
// 	if !c.config.Checkfile && len(checksum) == 0 {
// 		return false, utils.ErrCheckfile
// 	}

// 	configFile, err := os.Open(c.config.File)
// 	if err != nil {
// 		return false, err
// 	}

// 	defer configFile.Close()

// 	if checksum != "" {
// 		_, err = io.Copy(c.config.Algorithm, configFile)
// 		if err != nil {
// 			return false, err
// 		}

// 		return utils.Compare(checksum, c.config.Algorithm), nil
// 	}

// 	content := utils.SplitCheckfile(configFile)
// 	for _, line := range content {
// 		c.config.Algorithm.Reset()

// 		lineContent := strings.Split(line, "  ")

// 		f, err := os.Open(lineContent[1])
// 		if err != nil {
// 			return false, err
// 		}

// 		defer f.Close()

// 		_, err = io.Copy(c.config.Algorithm, f)
// 		if err != nil {
// 			return false, err
// 		}

// 		if !utils.Compare(lineContent[0], c.config.Algorithm) {
// 			return false, nil
// 		}
// 	}

// 	return true, nil
// }
