package gssc_test

import (
	"github.com/mh-cbon/gssc"
	"net/http"
	"time"
)

var port = ":8080"

// Example_main demonstrates usage of gssc package.
func Example_main() {
	s := &http.Server{
		Handler:      &ww{},
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			GetCertificate:     gssc.GetCertificate("example.org"),
		},
	}
	s.ListenAndServeTLS("", "")
}

type ww struct{}

func (s *ww) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}
