package langscanner

import (
	"fmt"
	"log"
	"strings"

	"github.com/dpastoor/wbi/internal/config"
)

var globalPythonPaths = []string{
	"/usr/bin/R",
	"/usr/local/bin/R",
	"/usr/lib/R",
}
var rootPythonDirs = []string{
	"/opt/R",
	"/usr/local/lib/R",
}

// GetPythonRootDirs returns the root directories for Python
func GetPythonRootDirs() []string {
	return rootPythonDirs
}

// GetPythonPaths returns the paths workbench will look for Python
// underneath the root directories with the format
// /root/{pythonVersion}/bin/python
func GetPythonPaths() []string {
	return globalPythonPaths
}

// ScanAndHandlePythonVersions scans for Python versions, handles result/errors and creates PythonConfig
func ScanAndHandlePythonVersions() (config.PythonConfig, error) {
	pythonVersions, err := ScanForPythonVersions()
	if err != nil {
		log.Fatal(err)
	}
	if len(pythonVersions) == 0 {
		// TODO: if no Python version found, send link to Python installation doc
		fmt.Println("no Pythons versions found at locations: \n", strings.Join(GetPythonRootDirs(), "\n"))
	} else {
		fmt.Println("found Python versions: ", strings.Join(pythonVersions, ", "))
	}

	var PythonConfig config.PythonConfig
	PythonConfig.Paths = pythonVersions

	return PythonConfig, err
}

// ScanForPythonVersions scans for Python versions in locations workbench will also look
func ScanForPythonVersions() ([]string, error) {
	// TODO actually do the scanning
	var versionsFound = []string{"/usr/bin/python"}

	return versionsFound, nil
}
