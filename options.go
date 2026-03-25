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
	Serve           bool
	ListenAddress   string
	Verbose         bool
}

func (o options) Validate() error {
	var missing []string

	if strings.TrimSpace(o.PrivateKeyPath) == "" {
		missing = append(missing, "-key")
	}
	if !o.Serve {
		if strings.TrimSpace(o.Issuer) == "" {
			missing = append(missing, "-iss")
		}
		if strings.TrimSpace(o.Subject) == "" {
			missing = append(missing, "-sub")
		}
		if strings.TrimSpace(o.Audience) == "" {
			missing = append(missing, "-aud")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
	}

	switch o.Algorithm {
	case "RS256", "PS256":
	default:
		return fmt.Errorf("unsupported alg %q, supported: RS256, PS256", o.Algorithm)
	}

	if o.Serve && strings.TrimSpace(o.ListenAddress) == "" {
		return fmt.Errorf("missing listen address for server mode")
	}

	if o.LifetimeSeconds <= 0 {
		return fmt.Errorf("lifetime must be greater than 0")
	}
	if o.GraceSeconds < 0 {
		return fmt.Errorf("grace must be 0 or greater")
	}

	return nil
}

func (o options) merge(overrides tokenRequest) options {
	merged := o

	if strings.TrimSpace(overrides.Algorithm) != "" {
		merged.Algorithm = overrides.Algorithm
	}
	if strings.TrimSpace(overrides.Issuer) != "" {
		merged.Issuer = overrides.Issuer
	}
	if strings.TrimSpace(overrides.Subject) != "" {
		merged.Subject = overrides.Subject
	}
	if strings.TrimSpace(overrides.Audience) != "" {
		merged.Audience = overrides.Audience
	}
	if strings.TrimSpace(overrides.Scope) != "" {
		merged.Scope = overrides.Scope
	}
	if overrides.LifetimeSeconds != nil {
		merged.LifetimeSeconds = *overrides.LifetimeSeconds
	}
	if overrides.IncludeNBF != nil {
		merged.IncludeNBF = *overrides.IncludeNBF
	}
	if overrides.IncludeIAT != nil {
		merged.IncludeIAT = *overrides.IncludeIAT
	}
	if overrides.GraceSeconds != nil {
		merged.GraceSeconds = *overrides.GraceSeconds
	}

	return merged
}

func (o options) ValidateRequest() error {
	var missing []string

	if strings.TrimSpace(o.Issuer) == "" {
		missing = append(missing, "iss")
	}
	if strings.TrimSpace(o.Subject) == "" {
		missing = append(missing, "sub")
	}
	if strings.TrimSpace(o.Audience) == "" {
		missing = append(missing, "aud")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required request fields: %s", strings.Join(missing, ", "))
	}

	switch o.Algorithm {
	case "RS256", "PS256":
	default:
		return fmt.Errorf("unsupported alg %q, supported: RS256, PS256", o.Algorithm)
	}

	if o.LifetimeSeconds <= 0 {
		return fmt.Errorf("lifetime must be greater than 0")
	}
	if o.GraceSeconds < 0 {
		return fmt.Errorf("grace must be 0 or greater")
	}

	return nil
}
