package main

import (
	"crypto/rsa"
	"fmt"
	"time"
)

type tokenService struct {
	key *rsa.PrivateKey
	kid string
}

func newTokenService(privateKeyPath string) (*tokenService, error) {
	key, err := loadRSAPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	kid, err := generateKid(&key.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate kid: %w", err)
	}

	return &tokenService{
		key: key,
		kid: kid,
	}, nil
}

func (s *tokenService) exchange(opts options, verbose bool) (string, error) {
	header, claims, err := buildTokenParts(opts, s.kid, time.Now().UTC())
	if err != nil {
		return "", err
	}

	assertion, encodedHeader, encodedClaims, err := signJWT(header, claims, s.key)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	if verbose {
		printAssertionDetails(assertion, s.kid, header, claims, encodedHeader, encodedClaims)
	}

	accessToken, err := exchangeToken(opts, assertion)
	if err != nil {
		return "", fmt.Errorf("failed to exchange token: %w", err)
	}

	return accessToken, nil
}
