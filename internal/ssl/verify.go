package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

// VerifySSLCertAndKeyMD5Match Verify the SSL cert and key are valid
func VerifySSLCertAndKeyMD5Match(certLocation string, keyLocation string) error {
	//Go has a method to do this, I opted not too, so that these commands can be rerun outside wbi
	certMD5, err := system.RunCommandAndCaptureOutput("openssl x509 -noout -modulus -in "+certLocation+" 2> /dev/null | openssl md5", false, 0, false)
	if err != nil {
		return fmt.Errorf("error checking md5 hash for certificate: %w", err)
	}
	keyMD5, err := system.RunCommandAndCaptureOutput("openssl rsa -noout -modulus -in "+keyLocation+" 2> /dev/null | openssl md5", false, 0, false)
	if err != nil {
		return fmt.Errorf("error checking md5 hash for key: %w", err)
	}

	if certMD5 != keyMD5 {
		return errors.New("Certificate and Key do not cryptographically match, " +
			"please check that your certificates are correct. You can run these commands to verify the certificates:\n" +
			"openssl x509 -noout -modulus -in " + certLocation + " 2> /dev/null | openssl md5\n" +
			"openssl rsa -noout -modulus -in " + keyLocation + "  2> /dev/null | openssl md5 \n\n" +
			"The returned MD5 hashes must match: %w")
	}
	return nil
}
func VerifySSLHostMatch(serverCert *x509.Certificate) (bool, error) {
	certHostMisMatch := false

	hostname, err := os.Hostname()

	if err != nil {
		return certHostMisMatch, fmt.Errorf("failed to retrieve hostname: %w", err)
	}

	fmt.Println("Detected Server Name: " + hostname)
	fmt.Println("Detected Certificate Primary DNS Name: " + serverCert.DNSNames[0])

	if !strings.Contains(serverCert.DNSNames[0], hostname) {
		certHostMisMatch = true
	}

	return certHostMisMatch, nil
}

func VerifyTrustedCertificate(serverCert *x509.Certificate, intermediateCA *x509.CertPool) (bool, error) {

	//Load system root CAs
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return false, fmt.Errorf("failed to load system root CAs: %w", err)
	}

	opts := x509.VerifyOptions{
		Roots:         rootCAs,
		Intermediates: intermediateCA,
	}

	_, err = serverCert.Verify(opts)
	if err != nil {
		fmt.Println("The server certificate is not trusted by the system:" + err.Error())
		return false, nil
	} else {
		fmt.Println("This certificate is trusted by system")
	}

	return true, nil
}

func DecodePemFiles(certInput []byte) tls.Certificate {
	var cert tls.Certificate
	certPEMBlock := certInput
	var certDERBlock *pem.Block
	for {
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		}
	}
	return cert
}

func ParseCertificateChain(certLocation string) (*x509.Certificate, *x509.CertPool, *x509.Certificate, error) {

	certData, err := os.ReadFile(certLocation)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read the certificate from\n"+
			"the provided path when verifying it's subject: %w", err)
	}

	certDataChain := DecodePemFiles(certData)
	intermediateCA := x509.NewCertPool()
	var serverCert *x509.Certificate
	var rootCert *x509.Certificate

	for _, cert := range certDataChain.Certificate {
		x509Cert, err := x509.ParseCertificate(cert)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error parsing certificates: %w", err)
		}
		if x509Cert.IsCA && x509Cert.Subject.CommonName != x509Cert.Issuer.CommonName {
			intermediateCA.AddCert(x509Cert)
			system.PrintAndLogInfo("This is an Intermediate Certificate:" + x509Cert.Subject.CommonName)
		} else if !x509Cert.IsCA {
			serverCert = x509Cert
			system.PrintAndLogInfo("This is the Server Certificate:" + x509Cert.Subject.String())
		} else if x509Cert.IsCA {
			rootCert = x509Cert
			system.PrintAndLogInfo("This is the Root Certificate:" + x509Cert.Subject.String())
		}
	}

	if serverCert == nil {
		return nil, nil, nil, fmt.Errorf("The server portion of your certificate is empty, or not in pem format")
	}
	if rootCert == nil {
		system.PrintAndLogInfo("The root portion of your certificate is empty," +
			" this is an odd, but not impossible configuration")
	}

	return serverCert, intermediateCA, rootCert, nil
}
