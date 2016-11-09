# gssc

[![GoDoc](https://godoc.org/github.com/mh-cbon/gssc?status.svg)](https://godoc.org/github.com/mh-cbon/gssc)

Easily starts an HTTPS server with self-signed certificates.

# Install

```sh
go get github.com/mh-cbon/gssc
glide install github.com/mh-cbon/gssc
```

# Usage

```go
package main

import (
  "net/http"
  "github.com/mh-cbon/gssc"
)

var port = ":8080"
var domain = "example.org"

// Example_main demonstrates usage of gssc package.
func main() {
  s := &http.Server{
    Handler: &ww{},
  	Addr: port,
    WriteTimeout: 15 * time.Second,
    ReadTimeout:  15 * time.Second,
  	TLSConfig: &tls.Config{
      InsecureSkipVerify: true,
      GetCertificate: gssc.GetCertificate(domain),
    },
  }
  s.ListenAndServeTLS("", "")
}

type ww struct{}
func (s *ww) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("This is an example server.\n"))
}
```

# Credits - read more

- sscg (https://github.com/mh-cbon/sscg)
- Go Authors (https://golang.org/src/crypto/tls/generate_cert.go)
- Stackoverflow (http://stackoverflow.com/questions/40412270/implement-tls-config-getcertificate-with-self-signed-certificates)
