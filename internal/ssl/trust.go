package ssl

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/sol-eng/wbi/internal/config"
	"os"
)

func TrustRootCertificate(rootCA *x509.Certificate, osType config.OperatingSystem) error {

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCA.Raw})
	err := os.WriteFile("cert.pem", pemCert, 0644)
	if err != nil {
		return fmt.Errorf("writing file to disk failed: ", err)
	}
	return err
}
