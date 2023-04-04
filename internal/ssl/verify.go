package ssl

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/sol-eng/wbi/internal/system"
	"os"
	"strings"
)

// VerifySSLCertAndKeyMD5Match Verify the SSL cert and key are valid
func VerifySSLCertAndKeyMD5Match(certLocation string, keyLocation string) error {
	//Go has a method to do this, I opted not too, so that these commands can be rerun outside wbi
	certMD5, err := system.RunCommandAndCaptureOutput("openssl x509 -noout -modulus -in "+certLocation+" 2> /dev/null | openssl md5", false, 0)
	if err != nil {
		return fmt.Errorf("error checking md5 hash for certificate: %w", err)
	}
	keyMD5, err := system.RunCommandAndCaptureOutput("openssl rsa -noout -modulus -in "+keyLocation+" 2> /dev/null | openssl md5", false, 0)
	if err != nil {
		return fmt.Errorf("error checking md5 hash for key: %w", err)
	}

	if certMD5 != keyMD5 {
		return fmt.Errorf("Certificate and Key do not cryptographically match, "+
			"please check that your certificates are correct. You can run these commands to verify the certificates:\n"+
			"openssl x509 -noout -modulus -in /path/to/your/certificate.crt 2> /dev/null | openssl md5\n"+
			"openssl rsa -noout -modulus -in /path/to/your/private.key  2> /dev/null | openssl md5 \n\n"+
			"The returned MD5 hashes must match: %w", err)
	}
	return nil
}
func VerifySSLHostMatch(certLocation string) (bool, error) {
	certHostMatch := false
	certData, err := os.ReadFile(certLocation)
	if err != nil {
		return certHostMatch, fmt.Errorf("failed to read the certificate from\n"+
			"the provided path when verifying it's subject: %w", err)
	}
	// Decode the certificate data
	block, _ := pem.Decode(certData)
	if block == nil {
		return certHostMatch, fmt.Errorf("failed to decode PEM block: %w", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return certHostMatch, fmt.Errorf("failed to parse certificate PEM data: %w", err)
	}

	hostname, err := os.Hostname()

	if err != nil {
		return certHostMatch, fmt.Errorf("failed to retrieve hostname: %w", err)
	}

	fmt.Println("Detected Server Name: " + hostname)
	fmt.Println("Detected Certificate Primary DNS Name: " + cert.DNSNames[0])

	if strings.Contains(cert.DNSNames[0], hostname) {
		certHostMatch = true
		return certHostMatch, nil
	}

	return certHostMatch, nil
}

func VerifyTrustedCertificate(certLocation string) error {

	certData, err := os.ReadFile(certLocation)
	if err != nil {
		return fmt.Errorf("failed to read the certificate from\n"+
			"the provided path when verifying it's subject: %w", err)
	}
	// Decode the certificate data
	block, _ := pem.Decode(certData)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block: %w", err)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate PEM data: %w", err)
	}

	// Load system root CAs
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return fmt.Errorf("failed to load system root CAs: %w", err)
	}

	// Verify certificate
	opts := x509.VerifyOptions{
		Roots: rootCAs,
	}
	if _, err := cert.Verify(opts); err != nil {
		fmt.Println("Certificate is not trusted by system:")
	} else {
		fmt.Println("Certificate is trusted by system")
	}

	return nil
}
