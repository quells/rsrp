package rsrp

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// RedirectRequest copies a request and changes the URL
func RedirectRequest(r *http.Request, newURL string) (newRequest *http.Request, err error) {
	newRequest, err = http.NewRequest(r.Method, newURL, r.Body)
	if err != nil {
		return
	}

	newRequest = newRequest.WithContext(r.Context())

	for k, vs := range r.Header {
		for _, v := range vs {
			newRequest.Header.Add(k, v)
		}
	}

	for _, c := range r.Cookies() {
		newRequest.AddCookie(c)
	}

	return
}

// RouteAll routes all requests based on the RouteRules provided
func RouteAll(rules []RouteRule) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, rule := range rules {
			if rule.Match.MatchString(r.URL.Path) {
				newURL := rule.RewriteLocation(r.URL.Path)

				newRequest, err := RedirectRequest(r, newURL)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
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

		http.Error(w, fmt.Sprintf("no route found for %s", r.URL.Path), http.StatusNotFound)
	}
}
