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

// Prompts user if they want to install Python and does the installation
func PromptAndInstallPython(osType string) ([]string, error) {
	installPythonChoice, err := PythonInstallPrompt()
	if err != nil {
		return []string{}, fmt.Errorf("issue selecting Python installation: %w", err)
	}
	if installPythonChoice {
		validPythonVersions, err := RetrieveValidPythonVersions()
		if err != nil {
			return []string{}, fmt.Errorf("issue retrieving Python versions: %w", err)
		}
		installPythonVersions, err := RSelectVersionsPrompt(validPythonVersions)
		if err != nil {
			return []string{}, fmt.Errorf("issue selecting Python versions: %w", err)
		}
		for _, pythonVersion := range installPythonVersions {
			err = DownloadAndInstallPython(pythonVersion, osType)
			if err != nil {
				return []string{}, fmt.Errorf("issue installing Python version: %w", err)
			}
		}
		return installPythonVersions, nil
	}
	return []string{}, nil
}

// ScanAndHandlePythonVersions scans for Python versions, handles result/errors and creates PythonConfig
func ScanAndHandlePythonVersions(osType string) ([]string, error) {
	pythonVersionsOrig, err := ScanForPythonVersions()
	if err != nil {
		return []string{}, fmt.Errorf("issue occured in scanning for Python versions: %w", err)
	}

	fmt.Println("\nFound Python versions: ", strings.Join(pythonVersionsOrig, ", "), "\n")

	if len(pythonVersionsOrig) == 0 {
		scanMessage := "no Python versions found at locations: \n" + strings.Join(GetPythonRootDirs(), "\n")
		fmt.Println(scanMessage)

		installedPythonVersion, err := PromptAndInstallPython(osType)
		if err != nil {
			return []string{}, fmt.Errorf("issue installing Python: %w", err)
		}
		if len(installedPythonVersion) == 0 {
			return []string{}, errors.New("no Python versions have been installed")
		}
	} else {
		anyOptLocations := []string{}
		for _, value := range pythonVersionsOrig {
			matched, err := regexp.MatchString(".*/opt.*", value)
			if err == nil && matched {
				anyOptLocations = append(anyOptLocations, value)
			}
		}
		if len(anyOptLocations) == 0 {
			fmt.Println("Posit recommends installing version of Python into the /opt directory to not conflict/rely on the system installed version of Python.")
		}
		_, err := PromptAndInstallPython(osType)
		if err != nil {
			return []string{}, fmt.Errorf("issue installing Python: %w", err)
		}
	}

	pythonVersions, err := ScanForPythonVersions()
	if err != nil {
		return []string{}, fmt.Errorf("issue occured in scanning for Python versions: %w", err)
	}

	fmt.Println("\nFound Python versions: ", strings.Join(pythonVersions, ", "))
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

	maybePython, err := exec.LookPath("python3")
	if err == nil {
		foundVersions = AppendIfMissing(foundVersions, maybePython)
	}

	return foundVersions, nil
}
