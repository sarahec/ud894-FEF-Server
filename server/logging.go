/*
Copyright 2016 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
