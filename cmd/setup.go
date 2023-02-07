package cmd

import (
	"fmt"

	"github.com/dpastoor/wbi/internal/authentication"
	"github.com/dpastoor/wbi/internal/config"
	"github.com/dpastoor/wbi/internal/jupyter"
	"github.com/dpastoor/wbi/internal/languages"
	"github.com/dpastoor/wbi/internal/license"
	"github.com/dpastoor/wbi/internal/ssl"
	"golang.org/x/exp/slices"

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

	//TODO: check if workbench installed

	// Determine OS
	// TODO switch back to function
	// osType, err := os.DetectOS()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	osType := "ubuntu22"

	// Languages
	selectedLanguages := languages.PromptAndRespond()
	languages.ScanAndHandleRVersions(&WBConfig.RConfig)
	if slices.Contains(selectedLanguages, "python") {
		languages.ScanAndHandlePythonVersions(&WBConfig.PythonConfig)
	}

	// Jupyter
	if len(WBConfig.PythonConfig.Paths) > 0 {
		jupyterChoice := jupyter.InstallPrompt()
		if jupyterChoice {
			jupyterPythonTarget, err := jupyter.KernelPrompt(&WBConfig.PythonConfig)
			if err != nil {
				jupyter.InstallJupyter(jupyterPythonTarget)
			}
		}
	}

	// SSL
	sslChoice := ssl.PromptSSL()
	if sslChoice {
		sslCertPath := ssl.PromptSSLFilePath()
		sslCertKeyPath := ssl.PromptSSLKeyFilePath()
		ssl.VerifySSLCertAndKey(sslCertPath, sslCertKeyPath)
	}

	// Authentication
	authChoice := authentication.ConvertAuthType(authentication.PromptAuthentication())
	WBConfig.AuthType = authChoice
	authentication.HandleAuthChoice(&WBConfig, osType)

	// Licensing
	licenseKey := license.PromptLicense()
	license.ActivateLicenseKey(licenseKey)

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
