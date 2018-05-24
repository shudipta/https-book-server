package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	organization    = "abc.bd"
	commonName      = "Book Server"
	duration        = 365
	caCertFilename  = "cert-generator/ca.crt"
	caKeyFilename   = "cert-generator/ca.key"
	srvCertFilename = "cert-generator/srv.crt"
	srvKeyFilename  = "cert-generator/srv.key"
	clCertFilename  = "cert-generator/cl.crt"
	clKeyFilename   = "cert-generator/cl.key"
	isClient        = false
)

var (
	caCert, srvCert, clCert *x509.Certificate
	priv, caPriv            *rsa.PrivateKey
)

func newCertificate(organization, commonName string, duration int, check int, addresses []string) *x509.Certificate {
	certificate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(duration) * time.Hour * 24),
		BasicConstraintsValid: true,
	}
	if check == 1 {
		certificate.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		certificate.IsCA = true
		//return &certificate
	} else {
		certificate.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		certificate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
		// if check == 2 {
		// 	certificate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		// } else {
		// 	certificate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		// }
	}
	//
	for i := 0; i < len(addresses); i++ {
		if ip := net.ParseIP(addresses[i]); ip != nil {
			certificate.IPAddresses = append(certificate.IPAddresses, ip)
		} else {
			certificate.DNSNames = append(certificate.DNSNames, addresses[i])
		}
	}

	return &certificate
}

func generate(certificate, parent x509.Certificate, certFilename, keyFilename string, isCA bool) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}
	var parKey *rsa.PrivateKey
	if isCA {
		caPriv = priv
		parKey = priv
	} else {
		parKey = caPriv
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	certificate.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatal("Failed to generate serial number:", err)
	}
	//serverCert, err := cert.NewSignedCert(cfgForServer, serverKey, caCert, caKey)
	//cert.NewSelfSignedCACert(cfg, caKey)
	//cert.NewSignedCert(cfgForServer, serverKey, caCert, caKey)
	derBytes, err := x509.CreateCertificate(rand.Reader, &certificate, &parent, &priv.PublicKey, parKey)
	if err != nil {
		log.Fatal("Failed to create certificate:", err)
	}

	certOut, err := os.Create(certFilename)
	defer certOut.Close()
	if err != nil {
		log.Fatal("Failed to open "+certFilename+" for writing:", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// permission 0600 means owner can read and write file
	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer keyOut.Close()
	if err != nil {
		log.Fatal("Failed to open key "+keyFilename+" for writing:", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	fmt.Println("Certificate generated successfully")
	fmt.Println("\tCertificate: ", certFilename)
	fmt.Println("\tPrivate Key: ", keyFilename)
}

func caCertPair() {
	// addresses := []string{"192.168.99.100"} //"localhost", "127.0.0.1",
	addresses := []string{"192.168.99.100", "localhost", "127.0.0.1"}
	caCert = newCertificate(organization, commonName, duration, 1, addresses)
	generate(*caCert, *caCert, caCertFilename, caKeyFilename, true)
}

func srvCertPair() {
	// addresses := []string{"192.168.99.100"} //"localhost", "127.0.0.1",
	addresses := []string{"192.168.99.100", "localhost", "127.0.0.1"}
	srvCert = newCertificate(organization, commonName, duration, 2, addresses)
	generate(*srvCert, *caCert, srvCertFilename, srvKeyFilename, false)
}

func clCertPair() {
	// addresses := []string{"192.168.99.100"} //"localhost", "127.0.0.1",
	addresses := []string{"192.168.99.100", "localhost", "127.0.0.1"}
	srvCert = newCertificate(organization, commonName, duration, 3, addresses)
	generate(*srvCert, *caCert, clCertFilename, clKeyFilename, false)
}

func main() {
	caCertPair()
	srvCertPair()
	clCertPair()
}
