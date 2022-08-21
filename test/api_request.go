package test

import "fmt"

type Request struct {
	Method  string
	URL     string
	Headers []string
	Params  map[string]interface{}
}

func (r *Request) GetURLWithQueryParams() string {
	u := r.URL
	if len(r.Params) > 0 {
		u = u + "?"
		for key, val := range r.Params {
			u = fmt.Sprintf("%s&%s=%v", u, key, val)
		}
	}
	return u
}
