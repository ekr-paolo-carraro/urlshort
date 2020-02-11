package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ekr-paolo-carraro/gophercises-2-urlshort/students/ekr-paolo-carraro/urlshort"
)

type defaultHandler struct{}

func (defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ciao!")
}

func main() {

	optArg := flag.String("source", "", "path to json or yaml source for short-urls")
	flag.Parse()

	defaultHandler := defaultHandler{}

	var err error
	var handf http.HandlerFunc

	if *optArg == "" {
		handf = defaultHandler.ServeHTTP

	} else {

		defaultShorts := map[string]string{
			"/e": "www.ecosia.com",
		}
		if strings.Contains(*optArg, ".json") {
			jsonSource, err := ioutil.ReadFile(*optArg)
			if err != nil {
				log.Fatalf("error on source: %v", err)
			}

			handf, err = urlshort.JSONHandler([]byte(jsonSource), urlshort.MapHandler(defaultShorts, defaultHandler))

		} else if strings.Contains(*optArg, ".yaml") {
			yamlSource, err := ioutil.ReadFile(*optArg)
			if err != nil {
				log.Fatalf("error on source: %v", err)
			}

			handf, err = urlshort.YAMLHandler([]byte(yamlSource), urlshort.MapHandler(defaultShorts, defaultHandler))
		}

		if err != nil {
			log.Fatalf("error %v", err)
		}
	}

	log.Fatal(http.ListenAndServe(":8080", handf))
}
