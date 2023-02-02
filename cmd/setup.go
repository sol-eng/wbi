package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dpastoor/wbi/internal/config"
	"github.com/dpastoor/wbi/internal/jupyter"
	"github.com/dpastoor/wbi/internal/langscanner"
	"github.com/dpastoor/wbi/internal/ssl"

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

	// Ask language question
	var qs = []*survey.Question{
		{
			Name: "languages",
			Prompt: &survey.MultiSelect{
				Message: "What languages will you use",
				Options: []string{"R", "python"},
				Default: []string{"R", "python"},
			},
		},
	}
	languageAnswers := struct {
		Languages []string `survey:"languages"`
	}{}
	err := survey.Ask(qs, &languageAnswers)
	fmt.Println("You just chose the languages: ", languageAnswers.Languages)
	WBConfig.RConfig, err = langscanner.ScanAndHandleRVersions()
	WBConfig.PythonConfig, err = langscanner.ScanAndHandlePythonVersions()

	// If python found -- setup jupyter or ask to setup jupyter or check
	if len(WBConfig.PythonConfig.Paths) > 0 {
		// Ask if jupyter should be installed
		jupyterInstallName := true
		jupyterInstallPrompt := &survey.Confirm{
			Message: "Would you like to install Jupyter?",
		}
		survey.AskOne(jupyterInstallPrompt, &jupyterInstallName)
		// If Jupyter, then do the install steps
		if jupyterInstallName {
			// Allow the user to select a version of Python to target
			jupyterPythonTarget := ""
			jupyterPythonPrompt := &survey.Select{
				Message: "Select a Python kernel to install Jupyter into:",
				Options: WBConfig.PythonConfig.Paths,
			}
			survey.AskOne(jupyterPythonPrompt, &jupyterPythonTarget)
			// Install Jupyter
			jupyterInstallError := jupyter.InstallJupyter(jupyterPythonTarget)
			if jupyterInstallError != nil {
				log.Fatal(jupyterInstallError)
			}
		}
	}

	// Handle SSL cert
	// * ask if want SSL
	sslCertName := false
	sslCertPrompt := &survey.Confirm{
		Message: "Would you like to use SSL?",
	}
	survey.AskOne(sslCertPrompt, &sslCertName)

	if sslCertName {
		// Ask for cert and key locations
		certLocationName := ""
		certLocationPrompt := &survey.Input{
			Message: "Filepath to SSL certificate:",
		}
		survey.AskOne(certLocationPrompt, &certLocationName)

		certKeyLocationName := ""
		certKeyLocationPrompt := &survey.Input{
			Message: "Filepath to SSL certificate key:",
		}
		survey.AskOne(certKeyLocationPrompt, &certKeyLocationName)
		// Check to make sure cert and key are valid
		certVerificationError := ssl.VerifySSLCertAndKey(certLocationName, certKeyLocationName)
		if certVerificationError != nil {
			log.Fatal(certVerificationError)
		}
	}

	// Handle authentication
	choosenAuthenticationName := ""
	choosenAuthenticationPrompt := &survey.Select{
		Message: "Choose an authentication method:",
		Options: []string{"SAML", "OIDC", "Active Directory/LDAP", "PAM", "Other"},
	}
	survey.AskOne(choosenAuthenticationPrompt, &choosenAuthenticationName)
	// TODO: Handle based on the choosen method

	// TODO: Handle license key

	return err
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
