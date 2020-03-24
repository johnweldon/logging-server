package main

import (
	"net/http"
	"strings"
)

func wpadHandler(w http.ResponseWriter, r *http.Request) {
	if matchesAny([]string{"wpad.dat", "proxy.pac"}, r.URL.Path) {
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
	}
}

func matchesAny(suffixes []string, path string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}
