package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/workbench"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type configCmd struct {
	cmd  *cobra.Command
	opts configOpts
}

type configOpts struct {
	certPath string
	keyPath  string
	url      string
	source   string
}

func newConfig(configOpts configOpts, item string) error {
	if item == "ssl" {
		err := workbench.WriteSSLConfig(configOpts.certPath, configOpts.keyPath, configOpts.url)
		if err != nil {
			return fmt.Errorf("failed to write SSL config for Workbench: %w", err)
		}
	} else if item == "repo" {
		err := workbench.WriteRepoConfig(configOpts.url, configOpts.source)
		if err != nil {
			return fmt.Errorf("failed to write repo config for Workbench: %w", err)
		}
	} else if item == "connect-url" {
		err := workbench.WriteConnectURLConfig(configOpts.url)
		if err != nil {
			return fmt.Errorf("failed to write Connect URL config for Workbench: %w", err)
		}
	} else {
		return fmt.Errorf("invalid item provided, please provide one of the following: ssl, repo, connect-url")
	}
	return nil
}

func setConfigOpts(configOpts *configOpts) {
	configOpts.certPath = viper.GetString("cert-path")
	configOpts.keyPath = viper.GetString("key-path")
	configOpts.url = viper.GetString("url")
	configOpts.source = viper.GetString("source")
}

func (opts *configOpts) Validate(args []string) error {
	// check args lengths
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided, please provide one argument")
	} else if len(args) > 1 {
		return fmt.Errorf("too many arguments provided, please provide only one argument")
	}

	// the cert-path flag is required for ssl
	if opts.certPath == "" && args[0] == "ssl" {
		return fmt.Errorf("the cert-path flag is required for ssl")
	}
	// the key-path flag is required for ssl
	if opts.keyPath == "" && args[0] == "ssl" {
		return fmt.Errorf("the key-path flag is required for ssl")
	}
	// the url flag is required for ssl
	if opts.url == "" && args[0] == "ssl" {
		return fmt.Errorf("the url flag is required for ssl")
	}

	// the cert-path flag is only valid for ssl
	if opts.certPath != "" && args[0] != "ssl" {
		return fmt.Errorf("the cert-path flag is only valid for ssl")
	}
	// the key-path flag is only valid for ssl
	if opts.keyPath != "" && args[0] != "ssl" {
		return fmt.Errorf("the key-path flag is only valid for ssl")
	}

	// the url flag is only valid for repo, connect-url and ssl
	if opts.url != "" && (args[0] != "repo" && args[0] != "connect-url" && args[0] != "ssl") {
		return fmt.Errorf("the url flag is only valid for repo, connect-url and url")
	}

	// the url flag is required for repo
	if opts.url == "" && args[0] == "repo" {
		return fmt.Errorf("the url flag is required for repo")
	}
	// the url flag is required for connect-url
	if opts.url == "" && args[0] == "connect-url" {
		return fmt.Errorf("the url flag is required for connect-url")
	}

	// the source flag is required for repo
	if opts.source == "" && args[0] == "repo" {
		return fmt.Errorf("the source flag is required for repo")
	}

	// the source flag is only valid for repo
	if opts.source != "" && args[0] != "repo" {
		return fmt.Errorf("the source flag is only valid for repo")
	}

	// the only source flags allow are cran and pypi
	if args[0] == "repo" && (opts.source != "cran" && opts.source != "pypi") {
		return fmt.Errorf("the source flag only allows cran and pypi")
	}

	return nil
}

func newConfigCmd() *configCmd {
	var configOpts configOpts

	root := &configCmd{opts: configOpts}

	// adding two spaces to have consistent formatting
	exampleText := []string{
		"To configure TLS/SSL:",
		"  wbi config ssl --cert-path [PATH-TO-CERTIFICATE-FILE] --key-path [PATH-TO-KEY-FILE] --url [SERVER-URL]",
		"",
		"To configure a default package repository:",
		"  wbi config repo --url [REPO-BASE-URL] --source cran",
		"  wbi config repo --url [REPO-BASE-URL] --source pypi",
		"",
		"To configure a default Posit Connect server:",
		"  wbi config connect-url --url [CONNECT-SERVER-URL]",
	}

	cmd := &cobra.Command{
		Use:     "config [item]",
		Short:   "Configure SSL, package repos, or a Connect server in Posit Workbench",
		Example: strings.Join(exampleText, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setConfigOpts(&root.opts)
			if err := root.opts.Validate(args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("config-opts")
			if err := newConfig(root.opts, strings.ToLower(args[0])); err != nil {
				return err
			}
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringP("cert-path", "c", "", "TLS/SSL certificate path")
	viper.BindPFlag("cert-path", cmd.Flags().Lookup("cert-path"))

	cmd.Flags().StringP("key-path", "k", "", "TLS/SSL key path")
	viper.BindPFlag("key-path", cmd.Flags().Lookup("key-path"))

	cmd.Flags().StringP("url", "u", "", "Package Manager, Connect or Server URL")
	viper.BindPFlag("url", cmd.Flags().Lookup("url"))

	cmd.Flags().StringP("source", "s", "", "Repository source (cran or pypi)")
	viper.BindPFlag("source", cmd.Flags().Lookup("source"))

	root.cmd = cmd
	return root
}
