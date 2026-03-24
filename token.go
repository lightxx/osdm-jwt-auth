package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
}

type jwtClaims struct {
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Aud   string `json:"aud"`
	Exp   int64  `json:"exp"`
	Scope string `json:"scope"`
	Nbf   *int64 `json:"nbf,omitempty"`
	Iat   *int64 `json:"iat,omitempty"`
	Jti   string `json:"jti"`
}

func buildTokenParts(opts options, kid string, now time.Time) (jwtHeader, jwtClaims, error) {
	nowUnix := now.Unix()
	header := jwtHeader{
		Alg: opts.Algorithm,
		Typ: "JWT",
		Kid: kid,
	}

	jti, err := newJTI()
	if err != nil {
		return jwtHeader{}, jwtClaims{}, fmt.Errorf("failed to generate jti: %w", err)
	}

	claims := jwtClaims{
		Iss:   opts.Issuer,
		Sub:   opts.Subject,
		Aud:   opts.Audience,
		Exp:   nowUnix + int64(opts.LifetimeSeconds),
		Scope: opts.Scope,
		Nbf:   optionalUnixTime(opts.IncludeNBF, nowUnix-int64(opts.GraceSeconds)),
		Iat:   optionalUnixTime(opts.IncludeIAT, nowUnix),
		Jti:   jti,
	}

	return header, claims, nil
}

func optionalUnixTime(include bool, value int64) *int64 {
	if !include {
		return nil
	}

	v := value
	return &v
}

func signJWT(header jwtHeader, claims jwtClaims, key *rsa.PrivateKey) (token string, encodedHeader string, encodedClaims string, err error) {
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", "", "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", "", "", err
	}

	encodedHeader = base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims = base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := encodedHeader + "." + encodedClaims

	sum := sha256.Sum256([]byte(signingInput))

	var sig []byte
	switch header.Alg {
	case "RS256":
		sig, err = rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:])
	case "PS256":
		sig, err = rsa.SignPSS(rand.Reader, key, crypto.SHA256, sum[:], &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
			Hash:       crypto.SHA256,
		})
	default:
		err = fmt.Errorf("unsupported alg %q", header.Alg)
	}
	if err != nil {
		return "", "", "", err
	}

	token = signingInput + "." + base64.RawURLEncoding.EncodeToString(sig)
	return token, encodedHeader, encodedClaims, nil
}

func newJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16]),
	), nil
}
