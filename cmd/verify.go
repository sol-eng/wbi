package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/connect"
	"github.com/sol-eng/wbi/internal/license"
	"github.com/sol-eng/wbi/internal/packagemanager"
	"github.com/sol-eng/wbi/internal/ssl"
	"github.com/sol-eng/wbi/internal/workbench"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type verifyCmd struct {
	cmd  *cobra.Command
	opts verifyOpts
}

type verifyOpts struct {
	url      string
	repo     string
	language string
	certPath string
	keyPath  string
}

func newVerify(verifyOpts verifyOpts, item string) error {

	if item == "packagemanager" {
		if verifyOpts.url != "" {
			// verify URL is valid first
			cleanPackageManagerURL, err := packagemanager.VerifyPackageManagerURL(verifyOpts.url)
			if err != nil {
				return fmt.Errorf("issue with reaching the Posit Package Manager URL: %w", err)
			}
			if verifyOpts.language != "" && verifyOpts.repo != "" {
				err = packagemanager.VerifyPackageManagerRepo(cleanPackageManagerURL, verifyOpts.repo, verifyOpts.language)
				if err != nil {
					return fmt.Errorf("issue with checking the Posit Package Manager repo: %w", err)
				}
			}
		} else {
			return fmt.Errorf("the url flag is required for packagemanager")
		}
	} else if item == "connect-url" {
		if verifyOpts.url != "" {
			_, err := connect.VerifyConnectURL(verifyOpts.url)
			if err != nil {
				return fmt.Errorf("issue with checking the Connect URL: %w", err)
			}
		} else {
			return fmt.Errorf("the url flag is required for connect")
		}
	} else if item == "workbench" {
		workbenchInstalled := workbench.VerifyWorkbench()
		if !workbenchInstalled {
			return fmt.Errorf("Workbench is not installed") //nolint:staticcheck
		}
	} else if item == "ssl" {
		err := ssl.VerifySSLCertAndKey(verifyOpts.certPath, verifyOpts.keyPath)
		if err != nil {
			return fmt.Errorf("issue with checking the SSL cert and key: %w", err)
		}
	} else if item == "license" {
		_, err := license.CheckLicenseActivation()
		if err != nil {
			return fmt.Errorf("issue in checking for license activation: %w", err)
		}
	}

	return nil
}

func setVerifyOpts(verifyOpts *verifyOpts) {
	verifyOpts.url = viper.GetString("urls")
	verifyOpts.repo = viper.GetString("repo")
	verifyOpts.language = viper.GetString("language")
	verifyOpts.certPath = viper.GetString("cert-path")
	verifyOpts.keyPath = viper.GetString("key-path")
}

func (opts *verifyOpts) Validate(args []string) error {
	// check args lengths
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided, please provide one argument")
	} else if len(args) > 1 {
		return fmt.Errorf("too many arguments provided, please provide only one argument")
	}

	// only the url flag is supported for packagemanager and connect-url
	if opts.url != "" && (args[0] != "packagemanager" && args[0] != "connect-url") {
		return fmt.Errorf("the url flag is only supported for packagemanager and connect-url")
	}
	// only the repo flag is supported for packagemanager
	if opts.repo != "" && args[0] != "packagemanager" {
		return fmt.Errorf("the repo flag is only supported for packagemanager")
	}
	// only the language flag is supported for packagemanager
	if opts.language != "" && args[0] != "packagemanager" {
		return fmt.Errorf("the language flag is only supported for packagemanager")
	}
	// only the cert-path flag is supported for ssl
	if opts.certPath != "" && args[0] != "ssl" {
		return fmt.Errorf("the cert-path flag is only supported for ssl")
	}
	// only the key-path flag is supported for ssl
	if opts.keyPath != "" && args[0] != "ssl" {
		return fmt.Errorf("the key-path flag is only supported for ssl")
	}

	// the url flag is required for packagemanager
	if opts.url == "" && args[0] == "packagemanager" {
		return fmt.Errorf("the url flag is required for packagemanager")
	}
	// if the repo flag is provided, the language flag must also be provided
	if opts.repo != "" && opts.language == "" {
		return fmt.Errorf("the language flag is required when the repo flag is provided")
	}
	// if the language flag is provided, the repo flag must also be provided
	if opts.language != "" && opts.repo == "" {
		return fmt.Errorf("the repo flag is required when the language flag is provided")
	}

	// the url flag is required for connect-url
	if opts.url == "" && args[0] == "connect-url" {
		return fmt.Errorf("the url flag is required for connect-url")
	}

	// the cert-path flag is required for ssl
	if opts.certPath == "" && args[0] == "ssl" {
		return fmt.Errorf("the cert-path flag is required for ssl")
	}
	// the key-path flag is required for ssl
	if opts.keyPath == "" && args[0] == "ssl" {
		return fmt.Errorf("the key-path flag is required for ssl")
	}

	return nil
}

func newVerifyCmd() *verifyCmd {
	var verifyOpts verifyOpts

	root := &verifyCmd{opts: verifyOpts}

	// adding two spaces to have consistent formatting
	exampleText := []string{
		"To verify a Package Manager URL is valid:",
		"  wbi verify packagemanager --url [URL]",
		"",
		"To verify a Package Manager URL is valid and the repo is valid:",
		"  wbi verify packagemanager --url [URL] --repo [REPO-NAME] --language r",
		"  wbi verify packagemanager --url [URL] --repo [REPO-NAME] --language python",
		"",
		"To verify a Connect URL is valid:",
		"  wbi verify connect-url --url [URL]",
		"",
		"To verify Workbench is installed:",
		"  wbi verify workbench",
		"",
		"To verify TLS/SSL cert and key are valid:",
		"  wbi verify ssl --cert-path [CERT-PATH] --key-path [KEY-PATH]",
		"",
		"To verify a license is activated:",
		"  wbi verify license",
	}

	cmd := &cobra.Command{
		Use:     "verify [item]",
		Short:   "Verify an item is installed, configured correctly and has network connectivity",
		Example: strings.Join(exampleText, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setVerifyOpts(&root.opts)
			if err := root.opts.Validate(args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("verify-opts")
			if err := newVerify(root.opts, strings.ToLower(args[0])); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP("url", "", "", "Package Manager or Connect base URL")
	viper.BindPFlag("urls", cmd.Flags().Lookup("url"))

	cmd.Flags().StringP("repo", "r", "", "Name of the Package Manager repository")
	viper.BindPFlag("repo", cmd.Flags().Lookup("repo"))

	cmd.Flags().StringP("language", "l", "", "The type of Package Manager repository, r or python")
	viper.BindPFlag("language", cmd.Flags().Lookup("language"))

	cmd.Flags().StringP("cert-path", "c", "", "TLS/SSL certificate path")
	viper.BindPFlag("cert-path", cmd.Flags().Lookup("cert-path"))

	cmd.Flags().StringP("key-path", "k", "", "TLS/SSL key path")
	viper.BindPFlag("key-path", cmd.Flags().Lookup("key-path"))

	root.cmd = cmd
	return root
}
