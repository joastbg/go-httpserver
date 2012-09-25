// Copyright 2012 Johan Astborg - joastbg@gmail.com

// HTTP server with URL dispatcher and delegate mapping

package main

import (
	"fmt"
	"net/http"
	"regexp"
)

// HTTPDelegate function declaration
type HTTPDelegate func(rw http.ResponseWriter, rq *http.Request)

// Configuration data
type webConfig struct {
	mm map[*regexp.Regexp] HTTPDelegate
}

// BEGIN: Delegates

func testDelegate() HTTPDelegate {
	return func(rw http.ResponseWriter, rq *http.Request) {
		fmt.Fprintf(rw, "<h1>Test</h1>" + rq.URL.Path)
	}
}

func helloDelegate() HTTPDelegate {
	return func(rw http.ResponseWriter, rq *http.Request) {
		fmt.Fprintf(rw, "<h1>Hello</h1>" + rq.URL.Path)
	}
}

func farDelegate() HTTPDelegate {
	return func(rw http.ResponseWriter, rq *http.Request) {
		fmt.Fprintf(rw, "<h1>Far</h1>" + rq.URL.Path)
	}
}

// END: Delegates

// Generate a new configuration
func NewConfig(mm map[*regexp.Regexp] HTTPDelegate) *webConfig {
	c := new(webConfig)
	c.mm = mm
	return c
}

// Handles HTTP requests using mapped delegates
func (w *webConfig) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	for k, v := range w.mm {
		if k.MatchString(rq.URL.Path) {
			v(rw, rq)
			return
		}
	}
	fmt.Fprintf(rw, "No HTTPDelegate attached to: " + rq.URL.Path)
}

func main() {

	a := make(map[*regexp.Regexp] HTTPDelegate)

	// Map regular expressions to delegates
	// Example URLs: 
	// 127.0.0.1:8080/test/1
	// 127.0.0.1:8080/hello/10
	// 127.0.0.1:8080/far.html

	reg1, _ := regexp.Compile("^/test/\\d/?$");
	reg2, _ := regexp.Compile("^/hello/\\d/?$");
	reg3, _ := regexp.Compile("far\\.html");

	a[reg1] = testDelegate()
	a[reg2] = helloDelegate()
	a[reg3] = farDelegate()

	// Start HTTP server on port 8000
	err := http.ListenAndServe(":8000", NewConfig(a))

	if err != nil {
		fmt.Printf("Server failed: ", err.Error())
	}
}
