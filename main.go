package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/urfave/negroni"
)

var (
	listen  = ":3883"
	verbose = false
	public  = "public"
	ignored = []func(*http.Request) bool{}
)

func main() {
	if p := os.Getenv("PORT"); p != "" {
		listen = ":" + p
	}
	if v := os.Getenv("VERBOSE"); v != "" {
		log.Printf("Verbose Logging enabled")
		verbose = true
	}
	if p := os.Getenv("PUBLIC_DIR"); p != "" {
		public = p
	}
	if g := os.Getenv("IGNORED_HOSTS"); g != "" {
		for _, hn := range strings.Split(g, ",") {
			ignored = append(ignored, IgnoreHost(strings.TrimSpace(hn)))
		}
	}

	n := negroni.New(
		negroni.NewRecovery(),
		newLogger(os.Stdout, verbose, ignored...),
		negroni.NewStatic(http.Dir(public)))

	s := &http.Server{
		Addr:           listen,
		Handler:        n,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	log.Printf("listening on %s\n", listen)
	log.Fatal(s.ListenAndServe())
}
