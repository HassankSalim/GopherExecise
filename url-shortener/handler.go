package urlshort

import (
	"net/http"
	"gopkg.in/yaml.v2"
	"fmt"
	"log"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		redirectUrl, exists := selectNewUrl(pathsToUrls, r.URL.Path)
		if exists {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
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
	// TODO: Implement this...
	yamlPaths, err := createMapFromYaml(yml)
	return MapHandler(yamlPaths, fallback), err
}

func selectNewUrl(pathsToUrls map[string]string, path string) (string, bool) {
	redirectPath, ok := pathsToUrls[path]
	return redirectPath, ok
}

func createMapFromYaml(yml []byte) (map[string]string, error) {
	var m []map[string]string
	res := make(map[string]string)
	err := yaml.Unmarshal([]byte(yml), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}
	for _, yamlMap := range m {
		k := yamlMap["path"]
		v := yamlMap["url"]
		res[k] = v
	}

	return res, nil
}