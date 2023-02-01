package langscanner

import (
	"errors"
	"os"
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

// GetRootDirs returns the root directories for R
func GetRootDirs() []string {
	return rootRDirs
}

// ScanForRVersions scans for R versions
func ScanForRVersions() ([]string, error) {
	foundVersions := []string{}
	for _, rPath := range globalRPaths {
		if _, err := os.Stat(rPath); err == nil {
			if !errors.Is(err, os.ErrNotExist) {
				return foundVersions, err
			}
			foundVersions = append(foundVersions, rPath)
		}

	}
	//return []string{"/opt/R/3.6.3/bin/R", "/opt/R/4.0.2/bin/R"}, nil
	return []string{}, nil
}
