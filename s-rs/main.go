package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

const (
	DEFAULT_PORT       = "8080"
	X_CF_FORWARDED_URL = "X-Cf-Forwarded-Url"
)

type SimpleRoundTripper struct {
	transport http.RoundTripper
}

func newSimpleRoundTripper() *SimpleRoundTripper {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSslValidation()},
	}
	return &SimpleRoundTripper{
		transport: transport,
	}
}

func (s *SimpleRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error

	response, err = s.transport.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	return response, err
}

func main() {
	http.Handle("/", newProxy())

	log.Fatal(http.ListenAndServe(":"+getPort(), nil))
}

func newProxy() http.Handler {
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			url, err := url.Parse(r.Header.Get(X_CF_FORWARDED_URL))
			if err != nil {
				log.Fatalln(err.Error())
			}

			r.URL = url
			r.Host = url.Host
		},
		Transport: newSimpleRoundTripper(),
	}
	return proxy
}

func getPort() string {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}
	return port
}

func skipSslValidation() bool {
	var skipSslValidation bool
	var err error
	if skipSslValidation, err = strconv.ParseBool(os.Getenv("SKIP_SSL_VALIDATION")); err != nil {
		skipSslValidation = true
	}
	return skipSslValidation
}
