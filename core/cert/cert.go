package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"
)

const (
	// ServerName  cert comman name
	ServerName = "crocodile"
)

// GenerateCert generate cert cert key file
func GenerateCert(pemkeydir string) error {
	pemkeydir = strings.TrimRight(pemkeydir, "/")
	_, err := os.Stat(pemkeydir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pemkeydir, 0755)
		if err != nil {
			return err
		}
	}

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		CommonName: ServerName,
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,               //KeyUsage 与 ExtKeyUsage 用来表明该证书是用来做服务器认证的
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, // 密钥扩展用途的序列
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 3650),
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 1024)

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk) //DER 格式
	certOut, err := os.Create(fmt.Sprintf("%s/cert.pem", pemkeydir))
	if err != nil {
		return fmt.Errorf("os.Create  failed: %w", err)
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return fmt.Errorf("pem.Encode failed: %w", err)
	}
	err = certOut.Close()
	if err != nil {
		return fmt.Errorf("certOut.Close failed: %w", err)
	}
	keyOut, err := os.Create(fmt.Sprintf("%s/key.pem", pemkeydir))
	if err != nil {
		return fmt.Errorf("os.Create  failed: %w", err)
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	if err != nil {
		return fmt.Errorf("pem.Encode failed: %w", err)
	}
	err = keyOut.Close()
	if err != nil {
		return fmt.Errorf("certOut.Close failed: %w", err)
	}
	return nil
}
