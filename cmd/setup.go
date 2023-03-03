package cmd

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/sol-eng/wbi/internal/authentication"
	"github.com/sol-eng/wbi/internal/config"
	"github.com/sol-eng/wbi/internal/connect"
	"github.com/sol-eng/wbi/internal/jupyter"
	"github.com/sol-eng/wbi/internal/languages"
	"github.com/sol-eng/wbi/internal/license"
	"github.com/sol-eng/wbi/internal/os"
	"github.com/sol-eng/wbi/internal/packagemanager"
	"github.com/sol-eng/wbi/internal/prodrivers"
	"github.com/sol-eng/wbi/internal/ssl"
	"github.com/sol-eng/wbi/internal/workbench"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type setupCmd struct {
	cmd  *cobra.Command
	opts setupOpts
}

type setupOpts struct {
}

func newSetup(setupOpts setupOpts) error {

	var WBConfig config.WBConfig

	fmt.Println("Welcome to the Workbench Installer!\n")

	// Check if running as root
	err := os.CheckIfRunningAsRoot()
	if err != nil {
		return err
	}

	// Determine OS and install pre-requisites
	osType, err := os.DetectOS()
	if err != nil {
		return err
	}
	err = os.InstallPrereqs(osType)
	if err != nil {
		return err
	}

	// Languages
	selectedLanguages, err := languages.PromptAndRespond()
	if err != nil {
		return fmt.Errorf("issue selecting languages: %w", err)
	}

	// R
	WBConfig.RConfig.Paths, err = languages.ScanAndHandleRVersions(osType)
	if err != nil {
		return fmt.Errorf("issue finding R locations: %w", err)
	}
	// remove any path that starts with /usr and only offer symlinks for those that don't (i.e. /opt directories)
	rPathsFiltered := languages.RemoveSystemRPaths(WBConfig.RConfig.Paths)
	// check if R and Rscript has already been symlinked
	rSymlinked := languages.CheckIfRSymlinkExists()
	rScriptSymlinked := languages.CheckIfRscriptSymlinkExists()
	if (len(rPathsFiltered) > 0) && !rSymlinked && !rScriptSymlinked {
		err = languages.PromptAndSetRSymlinks(rPathsFiltered)
		if err != nil {
			return fmt.Errorf("issue setting R symlinks: %w", err)
		}
	}
	if lo.Contains(selectedLanguages, "python") {
		WBConfig.PythonConfig.Paths, err = languages.ScanAndHandlePythonVersions(osType)
		if err != nil {
			return fmt.Errorf("issue finding Python locations: %w", err)
		}
		err = languages.PromptAndSetPythonPATH(WBConfig.PythonConfig.Paths)
		if err != nil {
			return fmt.Errorf("issue setting Python PATH: %w", err)
		}
	}

	workbenchInstalled := workbench.VerifyWorkbench()
	// If Workbench is not detected then prompt to install
	if !workbenchInstalled {
		installWorkbenchChoice, err := workbench.WorkbenchInstallPrompt()
		if err != nil {
			return fmt.Errorf("issue selecting Workbench installation: %w", err)
		}
		if installWorkbenchChoice {
			err := workbench.DownloadAndInstallWorkbench(osType)
			if err != nil {
				return fmt.Errorf("issue installing Workbench: %w", err)
			}
		} else {
			log.Fatal("Workbench installation is required to continue")
		}
	}

	// Licensing
	licenseActivationStatus, err := license.CheckLicenseActivation()
	if err != nil {
		return fmt.Errorf("issue in checking for license activation: %w", err)
	}

	if !licenseActivationStatus {
		licenseChoice, err := license.PromptLicenseChoice()
		if err != nil {
			return fmt.Errorf("issue in prompt for license activate choice: %w", err)
		}

		if licenseChoice {
			licenseKey, err := license.PromptLicense()
			if err != nil {
				return fmt.Errorf("issue entering license key: %w", err)
			}
			ActivateErr := license.ActivateLicenseKey(licenseKey)
			if ActivateErr != nil {
				return fmt.Errorf("issue activating license key: %w", ActivateErr)
			}
		}
	}

	// Jupyter
	if len(WBConfig.PythonConfig.Paths) > 0 {
		jupyterChoice, err := jupyter.InstallPrompt()
		if err != nil {
			return fmt.Errorf("issue selecting Jupyter: %w", err)
		}

		if jupyterChoice {
			jupyterPythonTarget, err := jupyter.KernelPrompt(&WBConfig.PythonConfig)
			if err != nil {
				return fmt.Errorf("issue selecting Python location for Jupyter: %w", err)
			}
			// the path to jupyter must be set in the config, not python
			pythonSubPath, err := languages.RemovePythonFromPath(jupyterPythonTarget)
			if err != nil {
				return fmt.Errorf("issue removing Python from path: %w", err)
			}
			jupyterPath := pythonSubPath + "/jupyter"
			WBConfig.PythonConfig.JupyterPath = jupyterPath

			if jupyterPythonTarget != "" {
				err := jupyter.InstallJupyter(jupyterPythonTarget)
				if err != nil {
					return fmt.Errorf("issue installing Jupyter: %w", err)
				}
			}
		}
	}

	// Pro Drivers
	proDriversExistingStatus, err := prodrivers.CheckExistingProDrivers()
	if err != nil {
		return fmt.Errorf("issue in checking for prior pro driver installation: %w", err)
	}
	if !proDriversExistingStatus {
		installProDriversChoice, err := prodrivers.ProDriversInstallPrompt()
		if err != nil {
			return fmt.Errorf("issue selecting Pro Drivers installation: %w", err)
		}
		if installProDriversChoice {
			err := prodrivers.DownloadAndInstallProDrivers(osType)
			if err != nil {
				return fmt.Errorf("issue installing Pro Drivers: %w", err)
			}
		}
	}

	// SSL
	WBConfig.SSLConfig.Using, err = ssl.PromptSSL()
	if err != nil {
		return fmt.Errorf("issue selecting if SSL is to be used: %w", err)
	}
	if WBConfig.SSLConfig.Using {
		WBConfig.SSLConfig.CertPath, err = ssl.PromptSSLFilePath()
		if err != nil {
			return fmt.Errorf("issue with the provided SSL cert path: %w", err)
		}
		WBConfig.SSLConfig.KeyPath, err = ssl.PromptSSLKeyFilePath()
		if err != nil {
			return fmt.Errorf("issue with the provided SSL cert key path: %w", err)
		}
		verifySSLCert := ssl.VerifySSLCertAndKey(WBConfig.SSLConfig.CertPath, WBConfig.SSLConfig.KeyPath)
		if verifySSLCert != nil {
			return fmt.Errorf("could not verify the SSL cert: %w", err)
		}
		fmt.Println("SSL successfully setup and verified")
	}

	// Authentication
	WBConfig.AuthConfig.Using, err = authentication.PromptAuth()
	if err != nil {
		return fmt.Errorf("issue selecting if Authentication is to be setup: %w", err)
	}
	if WBConfig.AuthConfig.Using {
		WBConfig.AuthConfig.AuthType, err = authentication.PromptAndConvertAuthType()
		if err != nil {
			return fmt.Errorf("issue entering and converting AuthType: %w", err)
		}
		AuthErr := authentication.HandleAuthChoice(&WBConfig, osType)
		if AuthErr != nil {
			return fmt.Errorf("issue handling authentication: %w", AuthErr)
		}
	}

	// Package Manager URL
	WBConfig.PackageManagerConfig.Using, err = packagemanager.PromptPackageManagerChoice()
	if err != nil {
		return fmt.Errorf("issue in prompt for Posit Package Manager choice: %w", err)
	}
	if WBConfig.PackageManagerConfig.Using {
		// prompt for base URL
		rawPackageManagerURL, err := packagemanager.PromptPackageManagerURL()
		if err != nil {
			return fmt.Errorf("issue entering Posit Package Manager URL: %w", err)
		}
		// verify URL is valid first
		cleanPackageManagerURL, err := packagemanager.VerifyPackageManagerURL(rawPackageManagerURL, false)
		if err != nil {
			return fmt.Errorf("issue with reaching the Posit Package Manager URL: %w", err)
		}
		// R repo
		WBConfig.PackageManagerConfig.RURL, err = packagemanager.PromptPackageManagerNameAndBuildURL(cleanPackageManagerURL, osType, "r")
		if err != nil {
			return fmt.Errorf("issue entering Posit Package Manager R repo and building URL: %w", err)
		}
		// Python repo
		WBConfig.PackageManagerConfig.PythonURL, err = packagemanager.PromptPackageManagerNameAndBuildURL(cleanPackageManagerURL, osType, "python")
		if err != nil {
			return fmt.Errorf("issue entering Posit Package Manager Python repo and building URL: %w", err)
		}
	} else {
		WBConfig.PackageManagerConfig.Using, err = packagemanager.PromptPublicPackageManagerChoice()
		if err != nil {
			return fmt.Errorf("issue in prompt for Posit Public Package Manager choice: %w", err)
		}
		if WBConfig.PackageManagerConfig.Using {
			// validate the public package manager URL can be reached
			_, err := packagemanager.VerifyPackageManagerURL("https://packagemanager.rstudio.com", true)
			if err != nil {
				return fmt.Errorf("issue with reaching the Posit Public Package Manager URL: %w", err)
			}

			WBConfig.PackageManagerConfig.RURL, err = packagemanager.BuildPublicPackageManagerFullURL(osType)
			if err != nil {
				return fmt.Errorf("issue with creating the full Posit Public Package Manager URL: %w", err)
			}
		}
	}

	// Connect URL
	WBConfig.ConnectConfig.Using, err = connect.PromptConnectChoice()
	if err != nil {
		return fmt.Errorf("issue in prompt for Connect URL choice: %w", err)
	}
	if WBConfig.ConnectConfig.Using {
		rawConnectURL, err := connect.PromptConnectURL()
		if err != nil {
			return fmt.Errorf("issue entering Connect URL: %w", err)
		}
		WBConfig.ConnectConfig.URL, err = connect.VerifyConnectURL(rawConnectURL)
		if err != nil {
			return fmt.Errorf("issue with checking the Connect URL: %w", err)
		}
	}

	// Config
	// first check if there are any config changes needed
	configChangesNeeded := WBConfig.DetectConfigChange()
	if configChangesNeeded {
		// print out the changes needed
		fmt.Println("\n=== The following configuration changes are needed:")
		WBConfig.ConfigStructToText()
		// ask the user if they want to write the config automatically
		configWriteChoice, err := config.PromptWriteConfig()
		if err != nil {
			return fmt.Errorf("issue selecting if config changes are to be written: %w", err)
		}
		if configWriteChoice {
			err := config.WriteConfig(WBConfig)
			if err != nil {
				return fmt.Errorf("issue writing config: %w", err)
			}
			err = workbench.RestartRStudioServerAndLauncher()
			if err != nil {
				return fmt.Errorf("issue restarting RStudio Server and Launcher: %w", err)
			}
		}
	} else {
		fmt.Println("\n=== No configuration changes are needed")
	}

	fmt.Println("\nThanks for using wbi! Please remember to make any needed manual configuration changes and restart RStudio Server and Launcher.")
	return nil
}

func setSetupOpts(setupOpts *setupOpts) {

}

func (opts *setupOpts) Validate() error {
	return nil
}

func newSetupCmd() *setupCmd {
	root := &setupCmd{opts: setupOpts{}}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setSetupOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("setup-opts")
			if err := newSetup(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
