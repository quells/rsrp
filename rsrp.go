package rsrp

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// RouteAll routes all requests based on the RouteRules provided
func RouteAll(rules []RouteRule) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, rule := range rules {
			if rule.Match.MatchString(r.URL.Path) {
				newPath := rule.Rewrite.Source.ReplaceAllString(r.URL.Path, rule.Rewrite.Destination)
				newURL := fmt.Sprintf("%s%s", rule.Destination, newPath)

				newRequest, err := http.NewRequest(r.Method, newURL, r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				for k, vs := range r.Header {
					for _, v := range vs {
						newRequest.Header.Add(k, v)
					}
				}

				client := &http.Client{}
				resp, err := client.Do(newRequest)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				for k, vs := range resp.Header {
					for _, v := range vs {
						w.Header().Add(k, v)
					}
				}
				w.WriteHeader(resp.StatusCode)

				w.Write(body)
				return
			}
		}
		w.Write([]byte(r.URL.Path))
	}
}
