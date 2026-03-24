package main

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var block *pem.Block
	for {
		block, data = pem.Decode(data)
		if block == nil {
			return nil, errors.New("no PEM block found")
		}

		switch block.Type {
		case "RSA PRIVATE KEY":
			key, parseErr := x509.ParsePKCS1PrivateKey(block.Bytes)
			if parseErr != nil {
				return nil, parseErr
			}
			return key, nil
		case "PRIVATE KEY":
			key, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
			if parseErr != nil {
				return nil, parseErr
			}

			rsaKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("PRIVATE KEY is not an RSA private key")
			}
			return rsaKey, nil
		}

		if len(data) == 0 {
			break
		}
	}

	return nil, errors.New("supported PEM block not found (expected RSA PRIVATE KEY or PRIVATE KEY)")
}

func generateKid(pub *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(der)
	return base64.RawURLEncoding.EncodeToString(hash[:]), nil
}
