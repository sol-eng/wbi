package cmd

import (
	"fmt"

	"github.com/dpastoor/wbi/internal/authentication"
	"github.com/dpastoor/wbi/internal/config"
	"github.com/dpastoor/wbi/internal/jupyter"
	"github.com/dpastoor/wbi/internal/languages"
	"github.com/dpastoor/wbi/internal/license"
	"github.com/dpastoor/wbi/internal/os"
	"github.com/dpastoor/wbi/internal/ssl"
	"github.com/dpastoor/wbi/internal/workbench"
	"github.com/samber/lo"

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

	// Workbench installation
	_, err := workbench.VerifyWorkbench()
	if err != nil {
		return err
	}

	// Determine OS
	osType, err := os.DetectOS()
	if err != nil {
		return err
	}

	// Languages
	selectedLanguages, err := languages.PromptAndRespond()
	if err != nil {
		return fmt.Errorf("issue selecting languages: %w", err)
	}

	WBConfig.RConfig.Paths, err = languages.ScanAndHandleRVersions()
	if err != nil {
		return fmt.Errorf("issue finding R locations: %w", err)
	}

	if lo.Contains(selectedLanguages, "python") {
		WBConfig.PythonConfig.Paths, err = languages.ScanAndHandlePythonVersions()
		if err != nil {
			return fmt.Errorf("issue finding Python locations: %w", err)
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

			if jupyterPythonTarget != "" {
				err := jupyter.InstallJupyter(jupyterPythonTarget)
				if err != nil {
					return fmt.Errorf("issue installing Jupyter: %w", err)
				}
			}
		}
	}

	// SSL
	useSSL, err := ssl.PromptSSL()
	if err != nil {
		return fmt.Errorf("issue selecting if SSL is to be used: %w", err)
	}
	if useSSL {
		sslCertPath, err := ssl.PromptSSLFilePath()
		if err != nil {
			return fmt.Errorf("issue with the provided SSL cert path: %w", err)
		}
		sslCertKeyPath, err := ssl.PromptSSLKeyFilePath()
		if err != nil {
			return fmt.Errorf("issue with the provided SSL cert key path: %w", err)
		}
		verifySSLCert := ssl.VerifySSLCertAndKey(sslCertPath, sslCertKeyPath)
		if verifySSLCert != nil {
			return fmt.Errorf("could not verify the SSL cert: %w", err)
		}
		fmt.Println("SSL successfully setup and verified")
	}

	// Authentication
	WBConfig.AuthType, err = authentication.PromptAndConvertAuthType()
	if err != nil {
		return fmt.Errorf("issue entering and converting AuthType: %w", err)
	}
	AuthErr := authentication.HandleAuthChoice(&WBConfig, osType)
	if AuthErr != nil {
		return fmt.Errorf("issue handling authentication: %w", AuthErr)
	}

	// Licensing
	licenseKey, err := license.PromptLicense()
	if err != nil {
		return fmt.Errorf("issue entering license key: %w", err)
	}
	ActivateErr := license.ActivateLicenseKey(licenseKey)
	if ActivateErr != nil {
		return fmt.Errorf("issue activating license key: %w", ActivateErr)
	}

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
