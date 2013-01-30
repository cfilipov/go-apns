// Copyright Ⓒ 2013 Cristian Filipov. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package apns

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// LoadPemFile reads a combined certificate+key pem file into memory.
func LoadPemFile(pemFile string) (cert tls.Certificate, err error) {
	pemBlock, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return
	}
	return LoadPem(pemBlock)
}

// LoadPem is similar to tls.X509KeyPair found in tls.go except that this 
// function reads all blocks from the same file.
func LoadPem(pemBlock []byte) (cert tls.Certificate, err error) {
	var block *pem.Block
	for {
		block, pemBlock = pem.Decode(pemBlock)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, block.Bytes)
		} else {
			break
		}
	}

	///////////////////////////////////////////////////////////////////////////
	// The rest of the code in this function is copied from the tls.X509KeyPair
	// implementation found at http://golang.org/src/pkg/crypto/tls/tls.go, 
	// with the exception of minor changes (no need to decode the next block).
	///////////////////////////////////////////////////////////////////////////

	if len(cert.Certificate) == 0 {
		err = errors.New("crypto/tls: failed to parse certificate PEM data")
		return
	}

	if block == nil {
		err = errors.New("crypto/tls: failed to parse key PEM data")
		return
	}

	// OpenSSL 0.9.8 generates PKCS#1 private keys by default, while
	// OpenSSL 1.0.0 generates PKCS#8 keys. We try both.
	var key *rsa.PrivateKey
	if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		var privKey interface{}
		if privKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			err = errors.New("crypto/tls: failed to parse key: " + err.Error())
			return
		}

		var ok bool
		if key, ok = privKey.(*rsa.PrivateKey); !ok {
			err = errors.New("crypto/tls: found non-RSA private key in PKCS#8 wrapping")
			return
		}
	}

	cert.PrivateKey = key

	// We don't need to parse the public key for TLS, but we so do anyway
	// to check that it looks sane and matches the private key.
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return
	}

	if x509Cert.PublicKeyAlgorithm != x509.RSA || x509Cert.PublicKey.(*rsa.PublicKey).N.Cmp(key.PublicKey.N) != 0 {
		err = errors.New("crypto/tls: private key does not match public key")
		return
	}

	return
}
