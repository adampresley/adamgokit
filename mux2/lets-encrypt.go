package mux2

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

/*
Provide a path to store certificates, and the domain they apply to
for Let's Encrypt certificate generation.
*/
type LetsEncryptConfig struct {
	CertPath string
	Domain   string
}

func startCertManager(config LetsEncryptConfig) *tls.Config {
	var (
		tlsConfig   *tls.Config
		certManager autocert.Manager
	)

	certManager = autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(config.CertPath),
		HostPolicy: autocert.HostWhitelist(config.Domain),
	}

	// Create a TLS config using the autocert manager
	tlsConfig = &tls.Config{
		GetCertificate: certManager.GetCertificate,
		NextProtos:     []string{"h2", "http/1.1"},
	}

	httpServer := &http.Server{
		Addr:    ":80",
		Handler: certManager.HTTPHandler(nil),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("error starting HTTP server on port 80 for Let's Encrypt: %v", err))
		}
	}()

	return tlsConfig
}
