package w3c_traceparent_generator

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
)

const defaultHeader = "traceparent"

// Config the plugin configuration.
type Config struct {
	HeaderName string
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName: defaultHeader,
	}
}

// W3CTraceParentGenerator Plugin Struct.
type W3CTraceParentGenerator struct {
	next       http.Handler
	headerName string
	name       string
}

// New creates a new W3CTraceparentGenerator instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.HeaderName) == 0 {
		return nil, fmt.Errorf("no header name provided")
	}
	return &W3CTraceParentGenerator{
		next:       next,
		name:       name,
		headerName: config.HeaderName,
	}, nil
}

// ServeHTTP defines the behaviour of the plugin
func (generator *W3CTraceParentGenerator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if len(req.Header.Get(generator.headerName)) > 0 {
		rw.Header().Add(generator.headerName, req.Header.Get(generator.headerName))
		generator.next.ServeHTTP(rw, req)
	} else {
		traceParent := "00-" + RandomHexaDecimalStringOfLength(32) + "-" + RandomHexaDecimalStringOfLength(16) + "-" + "00"
		req.Header.Add(generator.headerName, traceParent)
		rw.Header().Add(generator.headerName, traceParent)
		generator.next.ServeHTTP(rw, req)
	}
}

// RandomHexaDecimalStringOfLength returns a random hexadecimal string of length n.
func RandomHexaDecimalStringOfLength(n int) string {
	b := make([]byte, n/2)

	if _, err := rand.Read(b); err != nil {
		log.Printf("[ERROR] error generating hex value: %v", err)
		//Fallback to all zeroes
		for i := range b {
			b[i] = 0
		}
	}

	return hex.EncodeToString(b)[:n]
}
