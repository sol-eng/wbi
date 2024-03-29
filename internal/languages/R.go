package languages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/install"
	cmdlog "github.com/sol-eng/wbi/internal/logging"
	"github.com/sol-eng/wbi/internal/system"
)

var nonNumericRVersions = []string{
	"next", "devel",
}

type availableRVersions struct {
	RVersions []string `json:"r_versions"`
}

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

// PromptAndInstallR Prompts user if they want to install R and does the installation
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

		var installRVersions []string
		for {
			installRVersions, err = RSelectVersionsPrompt(validRVersions)
			if err != nil {
				return []string{}, fmt.Errorf("issue selecting R versions: %w", err)
			}
			if len(installRVersions) == 0 {
				system.PrintAndLogInfo(`No R versions selected. Please select at least one version to install.`)
			} else {
				break
			}
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
func ScanAndHandleRVersions(osType config.OperatingSystem) error {
	rVersionsOrig, err := ScanForRVersions()
	if err != nil {
		return fmt.Errorf("issue occured in scanning for R versions: %w", err)
	}
	system.PrintAndLogInfo("\nFound R versions:")
	system.PrintAndLogInfo(strings.Join(rVersionsOrig, "\n"))

	if len(rVersionsOrig) == 0 {
		scanMessage := "no R versions found at locations: \n" + strings.Join(GetRRootDirs(), "\n")
		system.PrintAndLogInfo(scanMessage)

		installedRVersion, err := PromptAndInstallR(osType)
		if err != nil {
			return fmt.Errorf("issue installing R: %w", err)
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
			system.PrintAndLogInfo("Posit recommends installing version of R into the /opt directory to not conflict/rely on the system installed version of R.")
		}
		installedRVersion, err := PromptAndInstallR(osType)
		if err != nil {
			return fmt.Errorf("issue installing R: %w", err)
		}
		if len(installedRVersion) > 0 {
			system.PrintAndLogInfo("\nThe following R versions have been installed:\n" + strings.Join(installedRVersion, "\n"))
		}
	}

	rVersions, err := ScanForRVersions()
	if err != nil {
		return fmt.Errorf("issue occured in scanning for R versions: %w", err)
	}

	err = CheckPromtAndSetRSymlinks(rVersions)
	if err != nil {
		return fmt.Errorf("issue setting R symlinks: %w", err)
	}

	system.PrintAndLogInfo("\nFound R versions:")
	system.PrintAndLogInfo(strings.Join(rVersions, "\n"))
	return nil
}

// AppendIfMissing Append to a string slice only if the string is not yet in the slice
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
	foundOptVersions := []string{}
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
					if rootPath == "/opt/R" {
						foundOptVersions = append(foundOptVersions, rpath)
					} else {
						foundVersions = append(foundVersions, rpath)
					}
				}
			}
		}

	}

	maybeR, err := exec.LookPath("R")
	if err == nil {
		foundVersions = AppendIfMissing(foundVersions, maybeR)
	}

	// sort /opt/R versions
	foundOptVersionsSortedPaths, err := sortOptRVersionPaths(foundOptVersions)
	if err != nil {
		return []string{}, fmt.Errorf("issue sorting /opt/R versions: %w", err)
	}

	finalVersionPathSorted := append(foundOptVersionsSortedPaths, foundVersions...)
	return finalVersionPathSorted, nil
}

func sortOptRVersionPaths(versionPaths []string) ([]string, error) {
	foundOptVersionsOnly := []string{}
	for _, optVersion := range versionPaths {
		i := strings.Index(optVersion, "R")
		j := strings.Index(optVersion, "bin")
		foundOptVersionsOnly = append(foundOptVersionsOnly, optVersion[i+2:j-1])
	}
	versions, err := ConvertStringSliceToVersionSlice(foundOptVersionsOnly)
	if err != nil {
		return []string{}, fmt.Errorf("issue converting string slice to version slice: %w", err)
	}
	sortedVersions := SortVersionsDesc(versions)
	foundOptVersionsSorted := ConvertVersionSliceToStringSlice(sortedVersions)
	foundOptVersionsSortedPaths := []string{}
	for _, optVersion := range foundOptVersionsSorted {
		foundOptVersionsSortedPaths = append(foundOptVersionsSortedPaths, "/opt/R/"+optVersion+"/bin/R")
	}
	return foundOptVersionsSortedPaths, nil
}

// RInstallPrompt Prompt users if they would like to install R versions
func RInstallPrompt() (bool, error) {
	name := true
	messageText := "Would you like to install version(s) of R?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the R install prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

func RetrieveValidRVersions() ([]string, error) {
	rVersionURL := "https://cdn.posit.co/r/versions.json"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, rVersionURL, nil)
	if err != nil {
		return []string{}, errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return []string{}, errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return []string{}, errors.New("error in HTTP status code")
	}

	var availVersions availableRVersions
	err = json.NewDecoder(res.Body).Decode(&availVersions)
	if err != nil {
		return []string{}, errors.New("error unmarshalling JSON data")
	}

	numericVersions, err := removeElements(availVersions.RVersions, nonNumericRVersions)
	if err != nil {
		return []string{}, errors.New("failed to remove non-numeric R versions")
	}
	versions, err := ConvertStringSliceToVersionSlice(numericVersions)
	if err != nil {
		return []string{}, fmt.Errorf("issue converting string slice to version slice: %w", err)
	}

	sortedVersions := SortVersionsDesc(versions)
	if err != nil {
		return []string{}, errors.New("failed to sort versions")
	}
	stringVersions := ConvertVersionSliceToStringSlice(sortedVersions)

	return stringVersions, nil

}

// RSelectVersionsPrompt Prompt asking users which R version(s) they would like to install
func RSelectVersionsPrompt(availableRVersions []string) ([]string, error) {
	messageText := "Which version(s) of R would you like to install?"
	var qs = []*survey.Question{
		{
			Name: "rversions",
			Prompt: &survey.MultiSelect{
				Message: messageText,
				Options: availableRVersions,
				Default: availableRVersions[0],
			},
		},
	}
	rVersionsAnswers := struct {
		RVersions []string `survey:"rversions"`
	}{}
	err := survey.Ask(qs, &rVersionsAnswers, survey.WithRemoveSelectAll(), survey.WithRemoveSelectNone())
	if err != nil {
		return []string{}, errors.New("there was an issue with the R versions selection prompt")
	}
	log.Info(messageText)
	log.Info(strings.Join(rVersionsAnswers.RVersions, ", "))
	return rVersionsAnswers.RVersions, nil
}

// DownloadAndInstallR Downloads the R installer, and installs R
func DownloadAndInstallR(rVersion string, osType config.OperatingSystem) error {
	// Create InstallerInfo with the proper information
	installerInfo, err := PopulateInstallerInfo("r", rVersion, osType)
	if err != nil {
		return fmt.Errorf("PopulateInstallerInfo: %w", err)
	}
	// Download installer
	installerPath, err := install.DownloadFile("R", installerInfo.URL, installerInfo.Name)
	if err != nil {
		return fmt.Errorf("DownloadR: %w", err)
	}
	// Install R
	err = install.InstallLanguage("r", installerPath, osType, rVersion)
	if err != nil {
		return fmt.Errorf("InstallLanguage: %w", err)
	}
	// save to command log
	installCommand, err := install.RetrieveInstallCommand(installerInfo.Name, osType)
	if err != nil {
		return fmt.Errorf("RetrieveInstallCommand: %w", err)
	}
	cmdlog.Info("curl -O " + installerInfo.URL)
	cmdlog.Info(installCommand)

	return nil
}

// PromptAndSetRSymlinks prompts user to set R symlinks
func PromptAndSetRSymlinks(rPaths []string) error {
	setRSymlinkChoice, err := RSymlinkPrompt()
	if err != nil {
		return fmt.Errorf("an issue occured during the selection of R symlink choice: %w", err)
	}
	if setRSymlinkChoice {
		RPathChoice, err := RLocationSymlinksPrompt(rPaths)
		if err != nil {
			return fmt.Errorf("issue selecting R binary to add symlinks: %w", err)
		}
		err = SetRSymlinks(RPathChoice)
		if err != nil {
			return fmt.Errorf("issue setting R symlinks: %w", err)
		}

		system.PrintAndLogInfo("\nThe selected R directory  " + RPathChoice + " has been successfully symlinked and will be available on the default system PATH.\n")
	}
	return nil
}

// RSymlinkPrompt asks users if they would like to set R symlinks
func RSymlinkPrompt() (bool, error) {
	name := true
	messageText := `Would you like to symlink a R version to make it available on PATH? This is recommended so Workbench can default to this version of R and users can type "R" in the terminal.`
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the symlink R prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

// RLocationSymlinksPrompt asks users which R binary they want to symlink
func RLocationSymlinksPrompt(rPaths []string) (string, error) {
	// Allow the user to select a version of R to target
	target := ""
	messageText := "Select a R binary to symlink:"
	prompt := &survey.Select{
		Message: messageText,
		Options: rPaths,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", errors.New("there was an issue with the R selection prompt for symlinking")
	}
	if target == "" {
		return target, errors.New("no R binary selected to be symlinked")
	}
	log.Info(messageText)
	log.Info(target)
	return target, nil
}

// SetRSymlinks sets the R symlinks (both R and Rscript)
func SetRSymlinks(rPath string) error {
	rCommand := "ln -s " + rPath + " /usr/local/bin/R"
	err := system.RunCommand(rCommand, true, 0, true)
	if err != nil {
		return fmt.Errorf("error setting R symlink with the command '%s': %w", rCommand, err)
	}
	rScriptCommand := "ln -s " + rPath + "script /usr/local/bin/Rscript"
	err = system.RunCommand(rScriptCommand, true, 0, true)
	if err != nil {
		return fmt.Errorf("error setting Rscript symlink with the command '%s': %w", rScriptCommand, err)
	}
	return nil
}

// RemoveSystemRPaths removes the system R paths from string slice
func RemoveSystemRPaths(rPaths []string) []string {
	var filteredRPaths []string
	for _, rPath := range rPaths {
		if !strings.Contains(rPath, "/usr/") {
			filteredRPaths = append(filteredRPaths, rPath)
		}
	}
	return filteredRPaths
}

// CheckIfRSymlinkExists checks if the R symlink exists
func CheckIfRSymlinkExists() bool {
	_, err := os.Stat("/usr/local/bin/R")
	if err != nil {
		return false
	}

	system.PrintAndLogInfo("\nAn existing R symlink has been detected (/usr/local/bin/R)")
	return true
}

// CheckIfRscriptSymlinkExists checks if the Rscript symlink exists
func CheckIfRscriptSymlinkExists() bool {
	_, err := os.Stat("/usr/local/bin/Rscript")
	if err != nil {
		return false
	}

	system.PrintAndLogInfo("\nAn existing Rscript symlink has been detected (/usr/local/bin/Rscript)")
	return true
}

func CheckAndSetRSymlinks(rPath string) error {
	// check if R and Rscript has already been symlinked
	rSymlinked := CheckIfRSymlinkExists()
	rScriptSymlinked := CheckIfRscriptSymlinkExists()
	if !rSymlinked && !rScriptSymlinked {
		err := SetRSymlinks(rPath)
		if err != nil {
			return fmt.Errorf("issue setting R symlinks: %w", err)
		}
	} else {
		system.PrintAndLogInfo("R and Rscript symlinks already exist, skipping symlink creation")
	}
	return nil
}

func CheckPromtAndSetRSymlinks(rPaths []string) error {
	// remove any path that starts with /usr and only offer symlinks for those that don't (i.e. /opt directories)
	rPathsFiltered := RemoveSystemRPaths(rPaths)
	// check if R and Rscript has already been symlinked
	rSymlinked := CheckIfRSymlinkExists()
	rScriptSymlinked := CheckIfRscriptSymlinkExists()
	if (len(rPathsFiltered) > 0) && !rSymlinked && !rScriptSymlinked {
		err := PromptAndSetRSymlinks(rPathsFiltered)
		if err != nil {
			return fmt.Errorf("issue setting R symlinks: %w", err)
		}
	}
	return nil
}

func ValidateRVersions(rVersions []string) error {
	availRVersions, err := RetrieveValidRVersions()
	if err != nil {
		return fmt.Errorf("error retrieving valid R versions: %w", err)
	}
	for _, rVersion := range rVersions {
		if !lo.Contains(availRVersions, rVersion) {
			return errors.New("version " + rVersion + " is not a valid R version")
		}
	}
	return nil
}
