package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type tokenRequest struct {
	Algorithm       string `json:"alg"`
	Issuer          string `json:"iss"`
	Subject         string `json:"sub"`
	Audience        string `json:"aud"`
	Scope           string `json:"scope"`
	LifetimeSeconds *int   `json:"lifetime"`
	IncludeNBF      *bool  `json:"nbf"`
	IncludeIAT      *bool  `json:"iat"`
	GraceSeconds    *int   `json:"grace"`
	Verbose         bool   `json:"verbose"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Error       string `json:"error,omitempty"`
}

func runServer(base options, svc *tokenService) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		handleTokenRequest(w, r, base, svc)
	})

	server := &http.Server{
		Addr:    base.ListenAddress,
		Handler: mux,
	}

	fmt.Fprintf(os.Stderr, "Listening on http://%s\n", base.ListenAddress)
	fmt.Fprintln(os.Stderr, "POST /token to fetch an access token")

	return server.ListenAndServe()
}

func handleTokenRequest(w http.ResponseWriter, r *http.Request, base options, svc *tokenService) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, tokenResponse{Error: "method not allowed"})
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	var req tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, tokenResponse{Error: fmt.Sprintf("invalid JSON body: %v", err)})
		return
	}

	effective := base.merge(req)
	if err := effective.ValidateRequest(); err != nil {
		writeJSON(w, http.StatusBadRequest, tokenResponse{Error: err.Error()})
		return
	}

	accessToken, err := svc.exchange(effective, base.Verbose || req.Verbose)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, tokenResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, tokenResponse{AccessToken: accessToken})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, errors.New("failed to encode response").Error(), http.StatusInternalServerError)
	}
}
