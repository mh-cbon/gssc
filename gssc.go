// Get self-signed certificate. Implements tls.Config.GetCertificate
// to provide an easy way to start an HTTPS server with self-signed certificate.
package gssc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/pkg/errors"
	"math/big"
	"net"
	"strings"
	"time"
)

// GetCertificte returns a function which generates a self-signed Certificate
// and implements tls.Config.GetCertificate.
//
// It takes a string(hosname) or a Certopts{} whith more spceific options.
//
// It panics if arg is not a string or a Certopts{}.
func GetCertificate(arg interface{}) func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	var opts Certopts
	var err error
	if host, ok := arg.(string); ok {
		opts = Certopts{
			RsaBits:   2048,
			IsCA:      true,
			Host:      host,
			ValidFrom: time.Now(),
		}
	} else if o, ok := arg.(Certopts); ok {
		opts = o
	} else {
		err = errors.New("Invalid arg type, must be string(hostname) or Certopt{...}")
	}

	cert, err := generate(opts)
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return cert, err
	}
}

// Certopts is a struct to define option to generate the certificate.
type Certopts struct {
	RsaBits   int
	Host      string
	IsCA      bool
	ValidFrom time.Time
	ValidFor  time.Duration
}

// generate a certificte for given options.
func generate(opts Certopts) (*tls.Certificate, error) {

	priv, err := rsa.GenerateKey(rand.Reader, opts.RsaBits)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}

	notAfter := opts.ValidFrom.Add(opts.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate serial number\n")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: opts.ValidFrom,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(opts.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if opts.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create certificate")
	}

	return &tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}, nil
}
