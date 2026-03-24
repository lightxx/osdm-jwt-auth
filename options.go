package main

import (
	"fmt"
	"strings"
)

type options struct {
	PrivateKeyPath  string
	Algorithm       string
	Issuer          string
	Subject         string
	Audience        string
	Scope           string
	LifetimeSeconds int
	IncludeNBF      bool
	IncludeIAT      bool
	GraceSeconds    int
	Verbose         bool
}

func (o options) Validate() error {
	var missing []string

	if strings.TrimSpace(o.PrivateKeyPath) == "" {
		missing = append(missing, "-key")
	}
	if strings.TrimSpace(o.Issuer) == "" {
		missing = append(missing, "-iss")
	}
	if strings.TrimSpace(o.Subject) == "" {
		missing = append(missing, "-sub")
	}
	if strings.TrimSpace(o.Audience) == "" {
		missing = append(missing, "-aud")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
	}

	switch o.Algorithm {
	case "RS256", "PS256":
		return nil
	default:
		return fmt.Errorf("unsupported alg %q, supported: RS256, PS256", o.Algorithm)
	}
}
