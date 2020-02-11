package urlshort

import (
	"encoding/json"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.RequestURI
		candidateRedirect, exist := pathsToUrls[path]
		if exist {
			if strings.HasPrefix(candidateRedirect, "http://") == false {
				candidateRedirect = "http://" + candidateRedirect
			}
			http.Redirect(w, r, candidateRedirect, http.StatusSeeOther)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	yItems := []shortItem{}
	err := yaml.Unmarshal(yml, &yItems)
	if err != nil {
		return nil, err
	}
	shortsData := map[string]string{}
	for _, yItem := range yItems {
		shortsData[yItem.Path] = yItem.Url
	}

	return MapHandler(shortsData, fallback), nil
}

//JSONHandler parse json input in the format
// { "path": "/some-path", "url":"www.some-url.com"}
//and return http.HandlerFunc o error about json parsing
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	jItems := []shortItem{}
	err := json.Unmarshal(jsonData, &jItems)
	if err != nil {
		return nil, err
	}

	shortsData := map[string]string{}
	for _, jItem := range jItems {
		shortsData[jItem.Path] = jItem.Url
	}

	return MapHandler(shortsData, fallback), nil
}

type shortItem struct {
	Path string
	Url  string
}
