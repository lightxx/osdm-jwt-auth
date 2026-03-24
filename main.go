package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		exitErr(err)
	}
}

func run() error {
	opts := parseOptions()

	if err := opts.Validate(); err != nil {
		return err
	}

	key, err := loadRSAPrivateKey(opts.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	kid, err := generateKid(&key.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to generate kid: %w", err)
	}

	header, claims, err := buildTokenParts(opts, kid, time.Now().UTC())
	if err != nil {
		return err
	}

	token, encodedHeader, encodedClaims, err := signJWT(header, claims, key)
	if err != nil {
		return fmt.Errorf("failed to sign token: %w", err)
	}

	if opts.Verbose {
		printAssertionDetails(token, kid, header, claims, encodedHeader, encodedClaims)
	}

	accessToken, err := exchangeToken(opts, token)
	if err != nil {
		return fmt.Errorf("failed to exchange token: %w", err)
	}

	fmt.Println(accessToken)

	return nil
}

func parseOptions() options {
	var opts options

	flag.StringVar(&opts.PrivateKeyPath, "key", "", "Path to RSA private key PEM file")
	flag.StringVar(&opts.Algorithm, "alg", "RS256", "Signing algorithm: RS256 or PS256")
	flag.StringVar(&opts.Issuer, "iss", "", "Issuer / client ID")
	flag.StringVar(&opts.Subject, "sub", "", "Subject / client ID")
	flag.StringVar(&opts.Audience, "aud", "", "Audience / token endpoint URL")
	flag.StringVar(&opts.Scope, "scope", "uic_osdm", "Scope")
	flag.IntVar(&opts.LifetimeSeconds, "lifetime", 300, "Token lifetime in seconds")
	flag.BoolVar(&opts.IncludeNBF, "nbf", true, "Include nbf")
	flag.BoolVar(&opts.IncludeIAT, "iat", true, "Include iat")
	flag.IntVar(&opts.GraceSeconds, "grace", 120, "Grace period in seconds for nbf")
	flag.BoolVar(&opts.Verbose, "verbose", false, "Print the client assertion details to stderr")
	flag.Parse()

	return opts
}

func printAssertionDetails(token, kid string, header jwtHeader, claims jwtClaims, encodedHeader, encodedClaims string) {
	fmt.Fprintln(os.Stderr, "Client assertion:")
	fmt.Fprintln(os.Stderr, token)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "kid:", kid)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Header:")
	fmt.Fprintln(os.Stderr, prettyJSON(header))
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Payload:")
	fmt.Fprintln(os.Stderr, prettyJSON(claims))
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Compact parts:")
	fmt.Fprintln(os.Stderr, "header =", encodedHeader)
	fmt.Fprintln(os.Stderr, "payload =", encodedClaims)
}
