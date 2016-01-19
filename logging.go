package main

import (
	"log"
	"net/http"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
}

func (w *LoggingResponseWriter) Write(data []byte) (int, error) {
	log.Printf("> %s\n\n", data)
	return w.ResponseWriter.Write(data)
}

func logWrapper(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &LoggingResponseWriter{w}
		log.Printf("< %s %s\n", r.Method, r.RequestURI) // TODO use httputil.DumpRequest?
		handler.ServeHTTP(writer, r)
	})
}
