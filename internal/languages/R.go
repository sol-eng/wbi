package languages

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dpastoor/wbi/internal/config"
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

// Detects if the path is an R directory
func isRDir(path string) (string, bool) {
	rpath := filepath.Join(path, "bin", "R")
	if _, err := os.Stat(rpath); err == nil {
		return rpath, true
	}
	return rpath, false
}

// Prompts user if they want to install R and does the installation
func PromptAndInstallR(osType config.OperatingSystem) ([]string, error) {
	installRChoice, err := RInstallPrompt()
	if err != nil {
		return []string{}, fmt.Errorf("issue selecting R installation: %w", err)
	}
	if installRChoice {
		validRVersions, err := RetrieveValidRVersions()
		if err != nil {
			return []string{}, fmt.Errorf("issue retrieving R versions: %w", err)
		}
		installRVersions, err := RSelectVersionsPrompt(validRVersions)
		if err != nil {
			return []string{}, fmt.Errorf("issue selecting R versions: %w", err)
		}
		for _, rVersion := range installRVersions {
			err = DownloadAndInstallR(rVersion, osType)
			if err != nil {
				return []string{}, fmt.Errorf("issue installing R version: %w", err)
			}
		}
		return installRVersions, nil
	}
	return []string{}, nil
}

// ScanAndHandleRVersions scans for R versions, handles result/errors and creates RConfig
func ScanAndHandleRVersions(osType config.OperatingSystem) ([]string, error) {
	rVersionsOrig, err := ScanForRVersions()
	if err != nil {
		return []string{}, fmt.Errorf("issue occured in scanning for R versions: %w", err)
	}
	fmt.Println("\nFound R versions: ", strings.Join(rVersionsOrig, ", "))

	if len(rVersionsOrig) == 0 {
		scanMessage := "no R versions found at locations: \n" + strings.Join(GetRRootDirs(), "\n")
		fmt.Println(scanMessage)

		installedRVersion, err := PromptAndInstallR(osType)
		if err != nil {
			return []string{}, fmt.Errorf("issue installing R: %w", err)
		}
		if len(installedRVersion) == 0 {
			log.Fatal("R must be installed to continue. Please install R and try again.")
		}
	} else {
		anyOptLocations := []string{}
		for _, value := range rVersionsOrig {
			matched, err := regexp.MatchString(".*/opt.*", value)
			if err == nil && matched {
				anyOptLocations = append(anyOptLocations, value)
			}
		}
		if len(anyOptLocations) == 0 {
			fmt.Println("Posit recommends installing version of R into the /opt directory to not conflict/rely on the system installed version of R.")
		}
		installedRVersion, err := PromptAndInstallR(osType)
		if err != nil {
			return []string{}, fmt.Errorf("issue installing R: %w", err)
		}

		fmt.Println("\nThe following R versions have been installed: ", strings.Join(installedRVersion, ", "))
	}

	rVersions, err := ScanForRVersions()
	if err != nil {
		return []string{}, fmt.Errorf("issue occured in scanning for R versions: %w", err)
	}

	fmt.Println("\nFound R versions: ", strings.Join(rVersions, ", "))
	return rVersions, nil
}

// Append to a string slice only if the string is not yet in the slice
func AppendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
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

	maybeR, err := exec.LookPath("R")
	if err == nil {
		foundVersions = AppendIfMissing(foundVersions, maybeR)
	}

	return foundVersions, nil
}
