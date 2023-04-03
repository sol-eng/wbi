package ssl

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/sol-eng/wbi/internal/system"
	"os"
)

// Verify the SSL cert and key are valid
func VerifySSLCertAndKey(certLocation string, keyLocation string) (bool, bool, error) {
	certMD5Match := false
	certHostMatch := false
	//Go has a method to do this, I opted not too, so that these commands can be rerun outside wbi
	certMD5, err := system.RunCommandAndCaptureOutput("openssl x509 -noout -modulus -in "+certLocation+" 2> /dev/null | openssl md5", false, 0)
	if err != nil {
		return certMD5Match, certHostMatch, errors.New("error checking md5 hash for certificate: ")
	}
	keyMD5, err := system.RunCommandAndCaptureOutput("openssl rsa -noout -modulus -in "+keyLocation+" 2> /dev/null | openssl md5", false, 0)
	if err != nil {
		return certMD5Match, certHostMatch, errors.New("error checking md5 hash for key: ")
	}

	if certMD5 != keyMD5 {
		return certMD5Match, certHostMatch, errors.New("Certificate and Key do not cryptographically match, " +
			"please check that your certificates are correct. You can run these commands to verify the certificates:\n" +
			"openssl x509 -noout -modulus -in /path/to/your/certificate.crt 2> /dev/null | openssl md5\n" +
			"openssl rsa -noout -modulus -in /path/to/your/private.key  2> /dev/null | openssl md5 \n\n" +
			"The returned MD5 hashes must match")
	} else {
		certMD5Match = true
	}
	certData, err := os.ReadFile(certLocation)
	if err != nil {
		return certMD5Match, certHostMatch, errors.New("failed to read the certificate from\n" +
			"the provided path when verifying it's subject")
	}
	// Decode the certificate data
	block, _ := pem.Decode(certData)
	if block == nil {
		return certMD5Match, certHostMatch, errors.New("failed to decode PEM block")
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return certMD5Match, certHostMatch, errors.New("failed to parse certificate PEM data")
	}

	hostname, err := os.Hostname()

	if err != nil {
		return certMD5Match, certHostMatch, errors.New("failed to retrieve hostname")
	}

	fmt.Println("Detected Server Name: " + hostname)
	fmt.Println("Detected Certificate Subject Name: " + cert.DNSNames[0])

	if hostname != cert.DNSNames[0] {
		return certMD5Match, certHostMatch, errors.New("The certificate subject and hostname of this system don't\n" +
			"match. This may be intentional because you're using a proxy or load balancer.\n" +
			"Otherwise this is likely the wrong certificate and you should retrieve the correct\n" +
			"certificate from your issuer.\n")
	} else {
		certHostMatch = true
	}

	return certMD5Match, certHostMatch, nil
}
