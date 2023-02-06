package external

import (
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}
