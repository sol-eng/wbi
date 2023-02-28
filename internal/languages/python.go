package languages

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/install"
	"github.com/sol-eng/wbi/internal/system"
)

var availablePythonVersions = []string{
	"3.11.2",
	"3.11.1",
	"3.11.0",
	"3.10.10",
	"3.10.9",
	"3.10.8",
	"3.10.7",
	"3.10.6",
	"3.10.5",
	"3.10.4",
	"3.10.3",
	"3.10.2",
	"3.10.1",
	"3.10.0",
	"3.9.16",
	"3.9.15",
	"3.9.14",
	"3.9.13",
	"3.9.12",
	"3.9.11",
	"3.9.10",
	"3.9.9",
	"3.9.8",
	"3.9.7",
	"3.9.6",
	"3.9.5",
	"3.9.4",
	"3.9.3",
	"3.9.2",
	"3.9.1",
	"3.9.0",
	"3.8.16",
	"3.8.15",
	"3.8.14",
	"3.8.13",
	"3.8.12",
	"3.8.11",
	"3.8.10",
	"3.8.9",
	"3.8.8",
	"3.8.7",
	"3.8.6",
	"3.8.5",
	"3.8.4",
	"3.8.3",
	"3.8.2",
	"3.8.1",
	"3.8.0",
	"3.7.16",
	"3.7.15",
	"3.7.14",
	"3.7.13",
	"3.7.12",
	"3.7.11",
	"3.7.10",
	"3.7.9",
	"3.7.8",
	"3.7.7",
	"3.7.6",
	"3.7.5",
	"3.7.4",
	"3.7.3",
	"3.7.2",
	"3.7.1",
	"3.7.0",
}

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

// PromptAndSetPythonPATH prompts user to set Python PATH
func PromptAndSetPythonPATH(pythonPaths []string) error {
	setPathPythonChoice, err := PythonPATHPrompt()
	if err != nil {
		return fmt.Errorf("issue selecting adding Python to PATH: %w", err)
	}
	if setPathPythonChoice {
		// Remove python/python3 from each path so other binaries will be available in path such as pip, jupyter, etc.
		pythonPathsBin, err := RemovePythonFromPathSlice(pythonPaths)
		if err != nil {
			return fmt.Errorf("issue removing python from slice of locations: %w", err)
		}
		pythonPathChoice, err := PythonLocationPATHPrompt(pythonPathsBin)
		if err != nil {
			return fmt.Errorf("issue selecting Python binary to add to PATH: %w", err)
		}
		err = system.AddToPATH(pythonPathChoice, "python")
		if err != nil {
			return fmt.Errorf("issue adding Python binary to PATH: %w", err)
		}

		fmt.Println("\nThe selected Python directory  " + pythonPathChoice + " has been successfully added to PATH in the /etc/profile.d/wbi_python.sh file.\n")
	}
	return nil
}

// PythonLocationPATHPrompt asks users which Python binary they want to add to PATH
func PythonLocationPATHPrompt(pythonPaths []string) (string, error) {
	// Allow the user to select a version of Python to target
	target := ""
	prompt := &survey.Select{
		Message: "Select a Python binary to add to PATH:",
		Options: pythonPaths,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", errors.New("there was an issue with the Python selection prompt for adding to PATH")
	}
	if target == "" {
		return target, errors.New("no Python binary selected to add to PATH")
	}
	return target, nil
}

// Prompts user if they want to install Python and does the installation
func PromptAndInstallPython(osType config.OperatingSystem) ([]string, error) {
	installPythonChoice, err := PythonInstallPrompt()
	if err != nil {
		return []string{}, fmt.Errorf("issue selecting Python installation: %w", err)
	}
	if installPythonChoice {
		validPythonVersions, err := RetrieveValidPythonVersions()
		if err != nil {
			return []string{}, fmt.Errorf("issue retrieving Python versions: %w", err)
		}
		installPythonVersions, err := PythonSelectVersionsPrompt(validPythonVersions)
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
func ScanAndHandlePythonVersions(osType config.OperatingSystem) ([]string, error) {
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

// Prompt users if they would like to install Python versions
func PythonInstallPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: "Would you like to install version(s) of Python?",
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Python install prompt")
	}
	return name, nil
}

// PythonPATHPrompt asks users if they would like to set Python PATH
func PythonPATHPrompt() (bool, error) {
	name := true
	prompt := &survey.Confirm{
		Message: `Would you like to add a Python version to PATH? This is recommended so users can type "python" and "pip" in the terminal to access this specified version of python and associated tools.`,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the Python set PATH prompt")
	}
	return name, nil
}

func RetrieveValidPythonVersions() ([]string, error) {
	// TODO make this dynamic based on https://cdn.posit.co/python/versions.json
	return availablePythonVersions, nil
}

// Prompt asking users which Python version(s) they would like to install
func PythonSelectVersionsPrompt(availablePythonVersions []string) ([]string, error) {
	var qs = []*survey.Question{
		{
			Name: "pythonVersions",
			Prompt: &survey.MultiSelect{
				Message: "Which version(s) of Python would you like to install?",
				Options: availablePythonVersions,
				Default: availablePythonVersions[0],
			},
		},
	}
	pythonVersionsAnswers := struct {
		PythonVersions []string `survey:"pythonVersions"`
	}{}
	err := survey.Ask(qs, &pythonVersionsAnswers)
	if err != nil {
		return []string{}, errors.New("there was an issue with the Python versions selection prompt")
	}
	if len(pythonVersionsAnswers.PythonVersions) == 0 {
		return []string{}, errors.New("at least one Python version must be selected")
	}
	return pythonVersionsAnswers.PythonVersions, nil
}

// Downloads the Python installer, and installs Python
func DownloadAndInstallPython(pythonVersion string, osType config.OperatingSystem) error {
	// Create InstallerInfoPython with the proper information
	installerInfo, err := PopulateInstallerInfo("python", pythonVersion, osType)
	if err != nil {
		return fmt.Errorf("PopulateInstallerInfoPython: %w", err)
	}
	// Download installer
	filepath, err := install.DownloadFile("Python", installerInfo.URL, installerInfo.Name)
	if err != nil {
		return fmt.Errorf("DownloadPython: %w", err)
	}
	// Install Python
	err = install.InstallLanguage("python", filepath, osType, pythonVersion)
	if err != nil {
		return fmt.Errorf("InstallLanguage: %w", err)
	}
	// Upgrade pip, setuptools, and wheel
	err = UpgradePythonTools(pythonVersion)
	if err != nil {
		return fmt.Errorf("UpgradePythonTools: %w", err)
	}

	return nil
}

func UpgradePythonTools(pythonVersion string) error {
	upgradeCommand := "/opt/python/" + pythonVersion + "/bin/pip install --upgrade --no-warn-script-location --disable-pip-version-check pip setuptools wheel"
	err := system.RunCommand(upgradeCommand)
	if err != nil {
		return fmt.Errorf("issue upgrading pip, setuptools and wheel for Python: %w", err)
	}

	successMessage := "\npip, setuptools and wheel have been upgraded for Python version " + pythonVersion + "\n"
	fmt.Println(successMessage)

	return nil
}

// RemovePythonFromPath removes python or python3 from the end of a path so the directory can be used
func RemovePythonFromPath(pythonPath string) (string, error) {
	if _, err := regexp.MatchString(".*/python.*", pythonPath); err == nil {
		i := strings.LastIndex(pythonPath, "/python")
		excludingLast := pythonPath[:i] + strings.Replace(pythonPath[i:], "/python", "", 1)
		return excludingLast, nil
	} else if _, err := regexp.MatchString(".*/python3.*", pythonPath); err == nil {
		i := strings.LastIndex(pythonPath, "/python3")
		excludingLast := pythonPath[:i] + strings.Replace(pythonPath[i:], "/python3", "", 1)
		return excludingLast, nil
	} else {
		return pythonPath, nil
	}
}

// RemovePythonFromPathSlice removes python or python3 from the end of a set of path strings in a slice so the directories can be used
func RemovePythonFromPathSlice(pythonPaths []string) ([]string, error) {
	var newPythonPaths []string
	for _, pythonPath := range pythonPaths {
		newPythonPath, err := RemovePythonFromPath(pythonPath)
		if err != nil {
			return []string{}, err
		}
		newPythonPaths = append(newPythonPaths, newPythonPath)
	}
	return newPythonPaths, nil
}
