package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func prettyJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}
