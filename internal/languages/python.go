package languages

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dpastoor/wbi/internal/config"
	"github.com/hairyhenderson/go-which"
	"github.com/samber/lo"
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
func ScanAndHandlePythonVersions(PythonConfig *config.PythonConfig) {
	pythonVersions, err := ScanForPythonVersions()
	if err != nil {
		log.Fatal(err)
	}
	if len(pythonVersions) == 0 {
		log.Fatal("no Python versions found at locations: \n", strings.Join(GetPythonRootDirs(), "\n"), "To install versions of Python, please follow the instructions outline here: https://docs.posit.co/resources/install-python/")
	}

	if !lo.Contains(pythonVersions, "/opt") {
		fmt.Println("Posit recommends installing version of Python into the /opt directory to not conflict/rely on the system installed version of Python. \n\nTo install versions of Python in this manner, please follow the instructions outline here: https://docs.posit.co/resources/install-python/")
	}

	fmt.Println("Found Python versions: ", strings.Join(pythonVersions, ", "))

	PythonConfig.Paths = pythonVersions
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

	maybePython := which.Which("python3")
	foundVersions = AppendIfMissing(foundVersions, maybePython)

	return foundVersions, nil
}
