package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	l := len(os.Args)
	fmt.Println("num of args is", l)
	privateKeyFile := os.Args[l-3]
	certificateFile := os.Args[l-2]
	caFile := os.Args[l-1]
	fmt.Println("=== private key", os.Args[l-3])
	fmt.Println("=== cert", os.Args[l-2])
	fmt.Println("=== ca", os.Args[l-1])

	// Load client cert
	cert, err := tls.LoadX509KeyPair(certificateFile, privateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		// InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	//resp, err := client.Get("https://192.168.99.100:30010")
	//resp, err := client.Get("https://localhost:8443/")
	resp, err := client.Get("https://https-book-server.default.svc:8443/")
	if err != nil {
		fmt.Println(err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("reading response error:", err)
	}
	fmt.Printf("%s\n", string(contents))
	for {
		fmt.Println(".")
		time.Sleep(time.Second*2)
	}
}
