package langscanner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpastoor/wbi/internal/config"
	"github.com/hairyhenderson/go-which"
)

var globalRPaths = []string{
	"/usr/bin/R",
	"/usr/local/bin/R",
	"/usr/lib/R",
}
var rootRDirs = []string{
	"/opt/R",
	"/usr/local/lib/R",
}

// GetRRootDirs returns the root directories for R
func GetRRootDirs() []string {
	return rootRDirs
}

// GetRPaths returns the paths workbench will look for R
// underneath the root directories with the format
// /root/{rversion}/bin/R
func GetRPaths() []string {
	return globalRPaths
}

func isRDir(path string) (string, bool) {
	rpath := filepath.Join(path, "bin", "R")
	if _, err := os.Stat(rpath); err == nil {
		return rpath, true
	}
	return rpath, false
}

// ScanAndHandleRVersions scans for R versions, handles result/errors and creates RConfig
func ScanAndHandleRVersions() (config.RConfig, error) {
	rVersions, err := ScanForRVersions()
	if err != nil {
		log.Fatal(err)
	}
	if len(rVersions) == 0 {
		// TODO: if no R version found, send link to R installation doc
		log.Fatal("no R versions found at locations: \n", strings.Join(GetRRootDirs(), "\n"))
	}

	fmt.Println("found R versions: ", strings.Join(rVersions, ", "))

	var RConfig config.RConfig
	RConfig.Paths = rVersions

	return RConfig, err
}

// ScanForRVersions scans for R versions in locations workbench will also look
func ScanForRVersions() ([]string, error) {
	foundVersions := []string{}
	// This is somewhat naive with respect to actually checking whether
	// this is _really_ a working version of R by launching it,
	// vs just matches the path to R
	for _, rPath := range GetRPaths() {
		if _, err := os.Stat(rPath); err == nil {
			foundVersions = append(foundVersions, rPath)
		}

	}
	for _, rootPath := range GetRRootDirs() {
		entries, err := os.ReadDir(rootPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return foundVersions, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				rpath, isR := isRDir(filepath.Join(rootPath, entry.Name()))
				if isR {
					foundVersions = append(foundVersions, rpath)
				}
			}
		}

	}
	maybeR := which.Which("R")
	if maybeR != "" {
		foundVersions = append(foundVersions, maybeR)
	}

	return foundVersions, nil
}
