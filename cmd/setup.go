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
	if lo.Contains(selectedLanguages, "r") {
		err := NewInstall(installOpts{versions: []string{}}, "r")
		if err != nil {
			return fmt.Errorf("issue installing R: %w", err)
		}
	} else {
		log.Fatal("R is required for Workbench")
	}

	if lo.Contains(selectedLanguages, "python") {
		err := NewInstall(installOpts{versions: []string{}}, "python")
		if err != nil {
			return fmt.Errorf("issue installing Python: %w", err)
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
			err := NewInstall(installOpts{versions: []string{}}, "workbench")
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
			err = NewActivate(activateOpts{key: licenseKey})
			if err != nil {
				return fmt.Errorf("issue activating license: %w", err)
			}
		}
	}

	// Jupyter
	if lo.Contains(selectedLanguages, "python") {
		jupyterChoice, err := jupyter.InstallPrompt()
		if err != nil {
			return fmt.Errorf("issue selecting Jupyter: %w", err)
		}
		if jupyterChoice {
			err := NewInstall(installOpts{versions: []string{}}, "jupyter")
			if err != nil {
				return fmt.Errorf("issue installing Jupyter: %w", err)
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
			err := NewInstall(installOpts{versions: []string{}}, "prodrivers")
			if err != nil {
				return fmt.Errorf("issue installing Pro Drivers: %w", err)
			}
		}
	}

	// SSL
	sslChoice, err := ssl.PromptSSL()
	if err != nil {
		return fmt.Errorf("issue selecting if SSL is to be used: %w", err)
	}
	if sslChoice {
		certPath, err := ssl.PromptSSLFilePath()
		if err != nil {
			return fmt.Errorf("issue with the provided SSL cert path: %w", err)
		}
		keyPath, err := ssl.PromptSSLKeyFilePath()
		if err != nil {
			return fmt.Errorf("issue with the provided SSL cert key path: %w", err)
		}
		err = NewVerify(verifyOpts{certPath: certPath, keyPath: keyPath}, "ssl")
		if err != nil {
			return fmt.Errorf("issue verifying SSL: %w", err)
		}
		err = NewConfig(configOpts{certPath: certPath, keyPath: keyPath}, "ssl")
		if err != nil {
			return fmt.Errorf("issue configuring SSL: %w", err)
		}

		fmt.Println("SSL successfully setup and verified")
	}

	// Authentication
	authChoice, err := authentication.PromptAuth()
	if err != nil {
		return fmt.Errorf("issue selecting if Authentication is to be setup: %w", err)
	}
	if authChoice {
		authType, err := authentication.PromptAndConvertAuthType()
		if err != nil {
			return fmt.Errorf("issue entering and converting AuthType: %w", err)
		}
		switch authType {
		case config.SAML:
			idpURL, err := authentication.PromptSAMLMetadataURL()
			if err != nil {
				return fmt.Errorf("issue entering SAML metadata URL: %w", err)
			}
			err = NewConfig(configOpts{authType: "saml", idpURL: idpURL}, "auth")
			if err != nil {
				return fmt.Errorf("issue configuring SAML: %w", err)
			}

			fmt.Println("Setting up SAML based authentication is a 2 step process. Step 1 was just completed by wbi writing the configuration file, however your SAML setup may require further configuration. \n\nTo complete step 2, you must configure your identify provider with Workbench following steps outlined here: https://docs.posit.co/ide/server-pro/authenticating_users/saml_sso.html#step-2.-configure-your-identity-provider-with-workbench")
		case config.OIDC:
			fmt.Println("Setting up OpenID Connect based authentication is a 2 step process. First configure your OpenID provider with the steps outlined here to complete step 1: https://docs.posit.co/ide/server-pro/authenticating_users/openid_connect_authentication.html#configuring-your-openid-provider \n\n As you register Workbench in the IdP, save the client-id and client-secret. Follow the next prompts to complete step 2.")

			clientID, err := authentication.PromptOIDCClientID()
			if err != nil {
				return fmt.Errorf("issue entering OIDC client ID: %w", err)
			}
			clientSecret, err := authentication.PromptOIDCClientSecret()
			if err != nil {
				return fmt.Errorf("issue entering OIDC client secret: %w", err)
			}
			idpURL, err := authentication.PromptOIDCIssuerURL()
			if err != nil {
				return fmt.Errorf("issue entering OIDC issuer URL: %w", err)
			}
			err = NewConfig(configOpts{authType: "oidc", idpURL: idpURL, clientID: clientID, clientSecret: clientSecret}, "auth")
			if err != nil {
				return fmt.Errorf("issue configuring OIDC: %w", err)
			}
		case config.LDAP:
			switch osType {
			case config.Ubuntu18, config.Ubuntu20, config.Ubuntu22:
				fmt.Println("Posit Workbench connects to LDAP via PAM. Please follow this article for more details on configuration: \nhttps://support.posit.co/hc/en-us/articles/360024137174-Integrating-Ubuntu-with-Active-Directory-for-RStudio-Workbench-RStudio-Server-Pro")
			case config.Redhat7, config.Redhat8:
				fmt.Println("Posit Workbench connects to LDAP via PAM. Please follow this article for more details on configuration: \nhttps://support.posit.co/hc/en-us/articles/360016587973-Integrating-RStudio-Workbench-RStudio-Server-Pro-with-Active-Directory-using-CentOS-RHEL")
			default:
				log.Fatal("Unsupported operating system")
			}
		case config.PAM:
			fmt.Println("PAM requires no additional configuration, however there are some username considerations and home directory provisioning steps that can be taken. To learn more please visit: https://docs.posit.co/ide/server-pro/authenticating_users/pam_authentication.html")
		case config.Other:
			fmt.Println("To learn about configuring your desired method of authentication please visit: https://docs.posit.co/ide/server-pro/authenticating_users/authenticating_users.html")
		}
	}

	// Package Manager URL
	packageManagerChoice, err := packagemanager.PromptPackageManagerChoice()
	if err != nil {
		return fmt.Errorf("issue in prompt for Posit Package Manager choice: %w", err)
	}
	if packageManagerChoice {
		// prompt for which languages to setup
		languageChoices, err := packagemanager.PromptLanguageRepos()
		if err != nil {
			return fmt.Errorf("issue in prompt for Posit Package Manager language choices: %w", err)
		}

		// prompt for base URL
		rawPackageManagerURL, err := packagemanager.PromptPackageManagerURL()
		if err != nil {
			return fmt.Errorf("issue entering Posit Package Manager URL: %w", err)
		}

		// r repo
		if lo.Contains(languageChoices, "r") {
			repoPackageManager, err := packagemanager.PromptPackageManagerRepo("r")
			if err != nil {
				return fmt.Errorf("issue entering Posit Package Manager repo name: %w", err)
			}
			cleanURL := packagemanager.CleanPackageManagerURL(rawPackageManagerURL)
			err = NewVerify(verifyOpts{url: cleanURL, repo: repoPackageManager, language: "r"}, "packagemanager")
			if err != nil {
				return fmt.Errorf("issue verifying Posit Package Manager URL: %w", err)
			}
			rPackageManagerURL, err := packagemanager.BuildPackagemanagerFullURL(cleanURL, repoPackageManager, osType, "r")
			if err != nil {
				return fmt.Errorf("issue building Posit Package Manager URL: %w", err)
			}
			err = NewConfig(configOpts{url: rPackageManagerURL, source: "cran"}, "repo")
			if err != nil {
				return fmt.Errorf("issue configuring Posit Package Manager CRAN URL: %w", err)
			}
		}

		// python repo
		if lo.Contains(languageChoices, "python") {
			repoPackageManager, err := packagemanager.PromptPackageManagerRepo("python")
			if err != nil {
				return fmt.Errorf("issue entering Posit Package Manager repo name: %w", err)
			}
			cleanURL := packagemanager.CleanPackageManagerURL(rawPackageManagerURL)
			err = NewVerify(verifyOpts{url: cleanURL, repo: repoPackageManager, language: "python"}, "packagemanager")
			if err != nil {
				return fmt.Errorf("issue verifying Posit Package Manager URL: %w", err)
			}
			pythonPackageManagerURL, err := packagemanager.BuildPackagemanagerFullURL(cleanURL, repoPackageManager, osType, "python")
			if err != nil {
				return fmt.Errorf("issue building Posit Package Manager URL: %w", err)
			}
			err = NewConfig(configOpts{url: pythonPackageManagerURL, source: "pypi"}, "repo")
			if err != nil {
				return fmt.Errorf("issue configuring Posit Package Manager PyPI URL: %w", err)
			}
		}
	} else {
		publicPackageManagerChoice, err := packagemanager.PromptPublicPackageManagerChoice()
		if err != nil {
			return fmt.Errorf("issue in prompt for Posit Public Package Manager choice: %w", err)
		}
		if publicPackageManagerChoice {
			publicPackageManagerURL := "https://packagemanager.rstudio.com"
			err = NewVerify(verifyOpts{url: publicPackageManagerURL}, "packagemanager")
			if err != nil {
				return fmt.Errorf("issue verifying Posit Public Package Manager URL: %w", err)
			}

			rPackageManagerURL, err := packagemanager.BuildPackagemanagerFullURL(publicPackageManagerURL, "cran", osType, "r")
			if err != nil {
				return fmt.Errorf("issue building Posit Public Package Manager URL: %w", err)
			}
			err = NewConfig(configOpts{url: rPackageManagerURL, source: "cran"}, "repo")
			if err != nil {
				return fmt.Errorf("issue configuring Posit Package Manager CRAN URL: %w", err)
			}
		}
	}

	// Connect URL
	connectChoice, err := connect.PromptConnectChoice()
	if err != nil {
		return fmt.Errorf("issue in prompt for Connect URL choice: %w", err)
	}
	if connectChoice {
		rawConnectURL, err := connect.PromptConnectURL()
		if err != nil {
			return fmt.Errorf("issue entering Connect URL: %w", err)
		}
		err = NewVerify(verifyOpts{url: rawConnectURL}, "connect-url")
		if err != nil {
			return fmt.Errorf("issue verifying Connect URL: %w", err)
		}
		cleanConnectURL := connect.CleanConnectURL(rawConnectURL)
		err = NewConfig(configOpts{url: cleanConnectURL}, "connect-url")
		if err != nil {
			return fmt.Errorf("issue configuring Connect URL: %w", err)
		}
	}

	fmt.Println("\n Restarting RStudio Server and Launcher...")

	err = workbench.RestartRStudioServerAndLauncher()
	if err != nil {
		return fmt.Errorf("issue restarting RStudio Server and Launcher: %w", err)
	}

	fmt.Println("\n Printing the status of RStudio Server and Launcher...")

	err = workbench.StatusRStudioServerAndLauncher()
	if err != nil {
		return fmt.Errorf("issue running status for RStudio Server and Launcher: %w", err)
	}

	fmt.Println("\nThanks for using wbi!")
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
		Short: "A series of interactive prompts to setup Posit Workbench",
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
