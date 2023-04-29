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
	certPath      string
	keyPath       string
	url           string
	source        string
	authType      string
	idpURL        string
	usernameClaim string
	clientID      string
	clientSecret  string
}

func newConfig(configOpts configOpts, item string) error {
	if item == "ssl" {
		err := workbench.WriteSSLConfig(configOpts.certPath, configOpts.keyPath, configOpts.url)
		if err != nil {
			return fmt.Errorf("failed to write SSL config for Workbench: %w", err)
		}
	} else if item == "auth" {
		if configOpts.authType == "saml" {
			err := workbench.WriteSAMLAuthConfig(configOpts.idpURL)
			if err != nil {
				return fmt.Errorf("failed to write SAML auth config for Workbench: %w", err)
			}
		} else if configOpts.authType == "oidc" {
			err := workbench.WriteOIDCAuthConfig(configOpts.idpURL, configOpts.usernameClaim, configOpts.clientID, configOpts.clientSecret)
			if err != nil {
				return fmt.Errorf("failed to write OIDC auth config for Workbench: %w", err)
			}
		} else {
			return fmt.Errorf("invalid auth type provided, please provide one of the following: saml, oidc")
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
		return fmt.Errorf("invalid item provided, please provide one of the following: ssl, auth, repo, connect-url")
	}
	return nil
}

func setConfigOpts(configOpts *configOpts) {
	configOpts.certPath = viper.GetString("cert-path")
	configOpts.keyPath = viper.GetString("key-path")
	configOpts.url = viper.GetString("url")
	configOpts.source = viper.GetString("source")
	configOpts.authType = viper.GetString("auth-type")
	configOpts.idpURL = viper.GetString("idp-url")
	configOpts.usernameClaim = viper.GetString("username-claim")
	configOpts.clientID = viper.GetString("client-id")
	configOpts.clientSecret = viper.GetString("client-secret")
}

func (opts *configOpts) Validate(args []string) error {
	// check args lengths
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided, please provide one argument")
	} else if len(args) > 1 {
		return fmt.Errorf("too many arguments provided, please provide only one argument")
	}

	// the auth-type flag is only valid for auth
	if opts.authType != "" && args[0] != "auth" {
		return fmt.Errorf("the auth-type flag is only valid for auth")
	}

	// the idp-url flag is only valid for argument auth and auth-type saml or oidc
	if opts.idpURL != "" && (args[0] != "auth" || (opts.authType != "saml" && opts.authType != "oidc")) {
		return fmt.Errorf("the idp-url flag is only valid with auth as an argument and a auth-type flag of saml or oidc")
	}

	// the username-claim flag is only valid for argument auth and auth-type oidc
	if opts.usernameClaim != "" && (args[0] != "auth" || opts.authType != "oidc") {
		return fmt.Errorf("the username-claim flag is only valid with auth as an argument and a auth-type flag of oidc")
	}
	// the client-id flag is only valid for argument auth and auth-type oidc
	if opts.clientID != "" && (args[0] != "auth" || opts.authType != "oidc") {
		return fmt.Errorf("the client-id flag is only valid with auth as an argument and a auth-type flag of oidc")
	}
	// the client-secret flag is only valid for argument auth and auth-type oidc
	if opts.clientSecret != "" && (args[0] != "auth" || opts.authType != "oidc") {
		return fmt.Errorf("the client-secret flag is only valid with auth as an argument and a auth-type flag of oidc")
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

	// the auth-type flag is required for auth
	if opts.authType == "" && args[0] == "auth" {
		return fmt.Errorf("the auth-type flag is required for auth")
	}

	// the idp-url flag is required for argument auth and auth-type saml
	if opts.idpURL == "" && args[0] == "auth" && (opts.authType == "saml" || opts.authType == "oidc") {
		return fmt.Errorf("the idp-url flag is required for argument auth and auth-type flag of saml or odic")
	}

	// the client-id flag is required for argument auth and auth-type oidc
	if opts.clientID == "" && args[0] == "auth" && opts.authType == "oidc" {
		return fmt.Errorf("the client-id flag is required for argument auth and auth-type flag of odic")
	}
	// the client-secret flag is required for argument auth and auth-type oidc
	if opts.clientSecret == "" && args[0] == "auth" && opts.authType == "oidc" {
		return fmt.Errorf("the client-secret flag is required for argument auth and auth-type flag of odic")
	}

	// the only valid authType flags are saml and oidc
	if args[0] == "auth" && (opts.authType != "saml" && opts.authType != "oidc") {
		return fmt.Errorf("the auth-type flag only allows saml and oidc")
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
		"To configure SAML Authentication:",
		"  wbi config auth --auth-type saml --idp-url [IDP-SAML-METADATA-URL]",
		"",
		"To configure OIDC Authentication:",
		"  wbi config auth --auth-type oidc --idp-url [IDP-OIC-ISSUER-URL] --client-id [CLIENT-ID] --client-secret [CLIENT-SECRET]",
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
		Short:   "Configure SSL, Authentication, package repos, or a Connect server in Posit Workbench",
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

	cmd.Flags().StringP("auth-type", "a", "", "Authentication type (saml or oidc)")
	viper.BindPFlag("auth-type", cmd.Flags().Lookup("auth-type"))

	cmd.Flags().StringP("idp-url", "i", "", "")
	viper.BindPFlag("idp-url", cmd.Flags().Lookup("idp-url"))

	cmd.Flags().StringP("username-claim", "", "", "IdP Metdata URL for SAML or OIDC")
	viper.BindPFlag("username-claim", cmd.Flags().Lookup("username-claim"))

	cmd.Flags().StringP("client-id", "", "", "OIDC Client ID")
	viper.BindPFlag("client-id", cmd.Flags().Lookup("client-id"))

	cmd.Flags().StringP("client-secret", "", "", "OIDC Client Secret")
	viper.BindPFlag("client-secret", cmd.Flags().Lookup("client-secret"))

	root.cmd = cmd
	return root
}
