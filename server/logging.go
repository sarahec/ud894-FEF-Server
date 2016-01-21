package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

const LogMethodAndURL = 1
const LogWholeRequest = 2

type LoggingResponseWriter struct {
	http.ResponseWriter
}

func (w *LoggingResponseWriter) Write(data []byte) (int, error) {
	log.Printf("> %s\n\n", data)
	return w.ResponseWriter.Write(data)
}

func logWrapper(level int, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &LoggingResponseWriter{w}
		switch level {
		case LogMethodAndURL:
			log.Printf("< %s %s\n", r.Method, r.RequestURI)
		case LogWholeRequest:
			dump, err := httputil.DumpRequest(r, true)
			if err == nil {
				fmt.Printf("< %s\n", dump)
			} else {
				log.Fatalf("** %s\n", err)
			}
		}

		// Keep it simple for now. Could use httputil.DumpRequest for
		// extra-verbose logging

		handler.ServeHTTP(writer, r)
	})
}
