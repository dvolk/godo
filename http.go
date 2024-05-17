// http related utilities

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// get the form or JSON data from the request,
// and ensure the expected fields are present
func PostGetData(w http.ResponseWriter, r *http.Request, expectedFields []string) (map[string]string, error) {
	// works with both json and web form requests
	out := map[string]string{}
	if r.Method != "POST" {
		return out, errors.New("wrong method")
	}
	if r.Header.Get("Content-Type") == "application/json" {
		// it's a json request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return out, errors.New("can't read body")
		}
		err = json.Unmarshal(body, &out)
		if err != nil {
			return out, errors.New("can't parse body as json")
		}
	} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			return out, errors.New("Can't parse form")
		}
		for key, values := range r.Form {
			fmt.Println(key, values)
			for _, str := range values {
				val, err := url.QueryUnescape(str)
				if err != nil {
					return out, errors.New(fmt.Sprintf("couldn't unescape %s", str))
				}
				out[key] = val
			}
		}
	}

	for _, expectedField := range expectedFields {
		if _, ok := out[expectedField]; !ok {
			return out, errors.New("missing field")
		}
	}
	return out, nil
}
