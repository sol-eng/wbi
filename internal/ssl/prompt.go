package ssl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/config"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/system"
)

// PromptSSL Prompt asking users if they wish to use SSL
func PromptSSL() (bool, error) {
	name := false
	messageText := "Would you like to use SSL?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the SSL prompt")
	}
	log.Info(messageText)
	log.Info(fmt.Sprintf("%v", name))
	return name, nil
}

func PromptAndVerifySSL(osType config.OperatingSystem) (string, string, error) {
	certPath, err := PromptSSLFilePath()
	if err != nil {
		return certPath, "", fmt.Errorf("issue with the provided SSL cert path: %w", err)
	}
	keyPath, err := PromptSSLKeyFilePath()
	if err != nil {
		return certPath, keyPath, fmt.Errorf("issue with the provided SSL cert key path: %w", err)
	}
	err = VerifySSLCertAndKeyMD5Match(certPath, keyPath)
	if err != nil {
		return certPath, keyPath, fmt.Errorf("could not verify the SSL cert: %w", err)
	}
	serverCert, intermediateCertPool, rootCert, err := ParseCertificateChain(certPath)
	if err != nil {
		return certPath, keyPath, fmt.Errorf("could not parse the certificate chain: %w", err)
	}
	var noRootOK bool
	if rootCert == nil {
		noRootOK, err = PromptRootCAMissing()
		if err != nil {
			return certPath, keyPath, fmt.Errorf("failure prompting for answer about RootCA missing: %w", err)
		}
		if !noRootOK {
			return certPath, keyPath, fmt.Errorf("no root CA Certificate in the certificate chain,"+
				" acquire the root and any necessary intermediate certificates and append them to the certificate file"+
				" in this order from top to bottom, server -> intermediate -> root. To restart the installer from this"+
				" step please use this command: 'wbi setup --step ssl': %w", err)
		} else {
			system.PrintAndLogInfo("Because your organization does not require intermediate and root" +
				" certificates to be presented to your corporate browsers, you will need to work with your SSL team" +
				" to get a copy of your organizations root and any intermediate certificates and add them to this" +
				" Linux machines trust store. These certificates are required to facilitate inter-server communication" +
				" between Workbench, Connect and Package Manager. An article on this topic can be found here:" +
				"https://support.posit.co/hc/en-us/articles/4416056988567-RStudio-Team-SSL-Considerations")
		}
	}

	certHostMisMatch, err := VerifySSLHostMatch(serverCert)

	if certHostMisMatch {
		proceed, err := PromptMisMatchedHostName()
		if err != nil {
			return certPath, keyPath, fmt.Errorf("hostname mismatch error: %w", err)
		}
		if !proceed {
			return certPath, keyPath, fmt.Errorf("hostname mismatch error, exit without proceeding: %w", err)
		}
	}

	verified, err := VerifyTrustedCertificate(serverCert, intermediateCertPool)
	if err != nil {
		return certPath, keyPath, fmt.Errorf("failure while trying to verify server trust of the SSL cert: %w", err)
	}
	if verified {
		system.PrintAndLogInfo("SSL successfully verified")
	} else {
		trust, err := PromptAddRootCAToTrustStore()
		if err != nil {
			return certPath, keyPath, fmt.Errorf("failure while prompting administrator for "+
				"certificate trust: %w", err)
		}
		if trust {
			err = TrustRootCertificate(rootCert, osType)
			if err != nil {
				return certPath, keyPath, fmt.Errorf("failure while trying to trust the SSL cert: %w", err)
			}
			//re-verify certificate via Bash, because SystemCertPool isn't refreshed in this context
			output, err := system.RunCommandAndCaptureOutput("openssl verify "+certPath, true, 1, false)
			if err != nil {
				return certPath, keyPath, fmt.Errorf("failure while trying to re-verify server trust of the SSL cert: %w", err)
			}
			if strings.Contains(output, "verification failed") {
				return certPath, keyPath, fmt.Errorf("failure while trying to re-verify server trust of the SSL cert: %w", err)
			}
			system.PrintAndLogInfo("SSL certificate successfully trusted")
		}
	}

	return certPath, keyPath, nil
}

// PromptSSLFilePath Prompt asking users for a filepath to their SSL cert
func PromptSSLFilePath() (string, error) {
	target := ""
	messageText := "Filepath to SSL certificate:"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", errors.New("there was an issue with the SSL cert path prompt")
	}
	log.Info(messageText)
	log.Info(target)
	return target, nil
}

// PromptSSLKeyFilePath Prompt asking users for a filepath to their SSL cert key
func PromptSSLKeyFilePath() (string, error) {
	target := ""
	messageText := "Filepath to SSL certificate key:"
	prompt := &survey.Input{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &target)
	if err != nil {
		return "", errors.New("there was an issue with the SSL cert key path prompt")
	}
	log.Info(messageText)
	log.Info(target)
	return target, nil
}

func PromptMisMatchedHostName() (bool, error) {
	name := false
	messageText := "The hostname of your server and the subject name in the certificate " +
		"don't match.\n This is common in configurations that include a load balancer " +
		"or a proxy.\n If you would like to exit the installer, resolve the certificate mismatch\n" +
		" and restart the installer at this step, you can run \"wbi setup --step ssl\" \n" +
		"Please confirm that you want to proceed with mismatched names above?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the SSL prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}

func PromptAddRootCAToTrustStore() (bool, error) {
	name := true
	messageText := "The certificate provided is not trusted by the system, this system level trust is usually required" +
		"\n to support connectivity between systems. Would you like to add this untrusted root certificate " +
		"\n to the system trust store?"
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the CA Trust prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}

func PromptRootCAMissing() (bool, error) {
	name := true
	messageText := "The certificate provided does not include a root Certificate Authority," +
		" otherwise known as a Root CA Certificate. Generally, Posit products require a chain of certificates from server to" +
		" your domains root certificate authority in order to present a valid HTTPS connection. In rare circumstances" +
		" this is not required. Does your corporate browser trust the server certificate that you're using without the" +
		" presence of a root certificate? If you are not sure, please select No."
	prompt := &survey.Confirm{
		Message: messageText,
	}
	err := survey.AskOne(prompt, &name)
	if err != nil {
		return false, errors.New("there was an issue with the missing Root Cert prompt")
	}
	log.Info(messageText)
	log.Info(name)
	return name, nil
}
