package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/negroni"
)

func main() {
	fmt.Fprintf(os.Stdout, "starting\n")
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), &headerLogger{}, negroni.NewStatic(http.Dir("public")))
	s := &http.Server{
		Addr:           ":3883",
		Handler:        n,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}
	log.Fatal(s.ListenAndServe())
}

type headerLogger struct {
	w io.Writer
}

func (l *headerLogger) logRequest(r *http.Request) {
	if l.w == nil {
		l.w = os.Stdout
	}
	fmt.Fprintf(l.w, " <====>\n")
	fmt.Fprintf(l.w, "Request: proto - %s; host - %s; remote - %s; method - %s; uri - %s\n", r.Proto, r.Host, r.RemoteAddr, r.Method, r.RequestURI)
	fmt.Fprintf(l.w, "Headers:\n")
	if err := r.Header.Write(l.w); err != nil {
		log.Printf("Error writing headers: %v", err)
	}
	fmt.Fprintf(l.w, "\n\n")
}

func (l *headerLogger) logResponse(w http.ResponseWriter) {
	if l.w == nil {
		l.w = os.Stdout
	}
	respHeaders := w.Header()
	fmt.Fprintf(l.w, "Response:\n")
	fmt.Fprintf(l.w, "Headers:\n")
	if err := respHeaders.Write(l.w); err != nil {
		log.Printf("Error writing headers: %v", err)
	}
	fmt.Fprintf(l.w, " <====>\n")
}

func (l *headerLogger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	l.logRequest(r)
	next(w, r)
	l.logResponse(w)
}
