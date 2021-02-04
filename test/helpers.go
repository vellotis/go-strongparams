package strongparamstest

import (
	"fmt"
	"net/http"
	"net/url"
)

func mockQueryValues(query string, args ...interface{}) url.Values {
	return mockRequestWithQuery(query, args...).URL.Query()
}

func mockRequestWithQuery(query string, args ...interface{}) *http.Request {
	return &http.Request{
		URL: &url.URL{
			RawQuery: fmt.Sprintf(query, args...),
		},
	}
}

func mockRequestWithPostForm(query string, args ...interface{}) *http.Request {
	return &http.Request{
		PostForm: mockQueryValues(query, args...),
	}
}