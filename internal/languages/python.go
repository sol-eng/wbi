package languages

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var globalPythonPaths = []string{
	"/usr/bin/python",
	"/usr/bin/Python",
	"/usr/local/bin/python",
	"/usr/local/bin/Python",
	"/usr/lib/python",
	"/usr/lib/Python",
}
var rootPythonDirs = []string{
	"/opt/python",
	"/opt/Python",
	"/usr/local/lib/python",
	"/usr/local/lib/Python",
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

func isPythonDir(path string) (string, bool) {
	pythonPath := filepath.Join(path, "bin", "python")
	if _, err := os.Stat(pythonPath); err == nil {
		return pythonPath, true
	}
	return pythonPath, false
}

// ScanAndHandlePythonVersions scans for Python versions, handles result/errors and creates PythonConfig
func ScanAndHandlePythonVersions() ([]string, error) {
	pythonVersions, err := ScanForPythonVersions()
	if err != nil {
		return []string{}, fmt.Errorf("issue occured in scanning for Python versions: %w", err)
	}
	if len(pythonVersions) == 0 {
		fmt.Println("To install versions of Python, please follow the instructions outline here: https://docs.posit.co/resources/install-python/")

		errorMessage := "no Python versions found at locations: \n" + strings.Join(GetPythonRootDirs(), "\n")
		return []string{}, errors.New(errorMessage)
	}

	var anyOptLocations = make([]string, 0)
	for _, value := range pythonVersions {
		matched, err := regexp.MatchString(".*/opt.*", value)
		if err == nil && matched {
			anyOptLocations = append(anyOptLocations, value)
		}
	}

	if len(anyOptLocations) == 0 {
		fmt.Println("Posit recommends installing version of Python into the /opt directory to not conflict/rely on the system installed version of Python. \nTo install versions of Python in this manner, please follow the instructions outline here: https://docs.posit.co/resources/install-python/")
	}

	fmt.Println("\nFound Python versions: ", strings.Join(pythonVersions, ", "), "\n")

	return pythonVersions, nil
}

// ScanForPythonVersions scans for Python versions in locations workbench will also look
func ScanForPythonVersions() ([]string, error) {
	foundVersions := []string{}
	// This is somewhat naive with respect to actually checking whether
	// this is _really_ a working version of Python by launching it,
	// vs just matches the path to Python
	for _, pyPath := range GetPythonPaths() {
		if _, err := os.Stat(pyPath); err == nil {
			foundVersions = append(foundVersions, pyPath)
		}

	}
	for _, rootPath := range GetPythonRootDirs() {
		entries, err := os.ReadDir(rootPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return foundVersions, err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				pythonPath, isPython := isPythonDir(filepath.Join(rootPath, entry.Name()))
				if isPython {
					foundVersions = append(foundVersions, pythonPath)
				}
			}
		}
	}

	maybePython, _ := exec.LookPath("python3")
	foundVersions = AppendIfMissing(foundVersions, maybePython)

	return foundVersions, nil
}
