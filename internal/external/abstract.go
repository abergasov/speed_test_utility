package external

import (
	"net/http"
	"net/url"
)

//go:generate mockgen -source=abstract.go -destination=abstract_external_mock.go -package=external
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}
