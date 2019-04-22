package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
)

func IgnoreHost(host string) func(*http.Request) bool {
	return func(r *http.Request) bool {
		hdr := r.Header.Get("X-Forwarded-Host")
		if hdr == "" {
			hdr = r.Header.Get("Host")
		}
		return hdr == host
	}
}

func newLogger(w io.Writer, verbose bool, excludeFilter ...func(*http.Request) bool) *logger {
	ll := log.New(w, "[http] ", log.LstdFlags|log.LUTC)
	return &logger{verbose: verbose, l: ll, o: w, exclude: excludeFilter}
}

type logger struct {
	verbose bool
	o       io.Writer
	l       *log.Logger
	exclude []func(*http.Request) bool
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if l.shouldLog(r) {
		l.requestLogger(r)
		rw, logResponse := l.responseLogger(r, w)
		defer logResponse()

		next.ServeHTTP(rw, r)
		return
	}
	next.ServeHTTP(w, r)
}

func (l *logger) shouldLog(r *http.Request) bool {
	for _, fn := range l.exclude {
		if fn(r) {
			return false
		}
	}
	return true
}

func (l *logger) requestLogger(r *http.Request) {
	l.l.Printf("(client %s) %s %s %s [%s]", queryIP(r), r.Host, r.Method, r.URL.Path, r.UserAgent())
	if l.verbose {
		fmt.Fprint(l.o, "|REQUEST|\n")
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			l.l.Printf("error dumping request: %v", err)
			return
		}
		switch ct := r.Header.Get("Content-Type"); {
		case ct == "":
			l.l.Print("missing content-type")
		case strings.Contains(ct, "text/"), ct == "application/json":
			if _, err = l.o.Write(b); err != nil {
				l.l.Printf("error writing dumped request: %v", err)
				return
			}
		default:
			l.l.Printf("filter dump content-type %q:\n", ct)
			f := filter(l.o)
			if _, err = f.Write(b); err != nil {
				l.l.Printf("error writing dumped request: %v", err)
				return
			}
			l.l.Print("\n")
		}
	}
}

type binaryFilter struct{ w io.Writer }

func (f *binaryFilter) Write(p []byte) (n int, err error) {
	for ix, b := range p {
		if b < ' ' {
			p[ix] = '?'
		}
		if b > '~' {
			p[ix] = '?'
		}
	}
	return f.w.Write(p)
}

func filter(w io.Writer) io.Writer {
	return &binaryFilter{w: w}
}

func (l *logger) responseLogger(r *http.Request, w http.ResponseWriter) (http.ResponseWriter, func()) {
	rr := httptest.NewRecorder()
	return rr, func() {
		for k, v := range rr.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rr.Code)

		out := []io.Writer{w}

		if l.verbose {
			fmt.Fprint(l.o, "|RESPONSE|\n")
			if err := rr.HeaderMap.Write(l.o); err != nil {
				l.l.Printf("error dumping response headers: %v", err)
			}
			switch ct := rr.HeaderMap.Get("Content-Type"); {
			case ct == "":
				// ignore no content type
			case strings.Contains(ct, "text/"):
				out = append(out, l.o)
			default:
				l.l.Printf("not logging content type %q", ct)
			}
		}

		if _, err := rr.Body.WriteTo(io.MultiWriter(out...)); err != nil {
			l.l.Printf("error sending response: %v", err)
		}

		if l.verbose {
			fmt.Fprint(l.o, "\n")
		}

		l.l.Printf("(client %s) %d %s", queryIP(r), rr.Code, http.StatusText(rr.Code))
	}
}
