package testutils

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"
)

// APIMocker is a mock for the http.Client interface.
type APIMocker struct {
	sleepDuration *time.Duration
	Getter        func(url string) (resp *http.Response, err error)
	Poster        func(url string, data url.Values) (resp *http.Response, err error)
	Doer          func(req *http.Request) (*http.Response, error)
}

// Get is a mock for the http.Client.Get method.
func (t *APIMocker) Get(urlAddr string) (resp *http.Response, err error) {
	if t.sleepDuration != nil {
		time.Sleep(*t.sleepDuration)
	}
	return t.Getter(urlAddr)
}

// PostForm is a mock for the http.Client.PostForm method.
func (t *APIMocker) PostForm(urlAddr string, data url.Values) (resp *http.Response, err error) {
	if t.sleepDuration != nil {
		time.Sleep(*t.sleepDuration)
	}
	return t.Poster(urlAddr, data)
}

// Do is a mock for the http.Client.Do method.
func (t *APIMocker) Do(req *http.Request) (*http.Response, error) {
	if t.sleepDuration != nil {
		time.Sleep(*t.sleepDuration)
	}
	return t.Doer(req)
}

// Option is a function that can be passed to New to modify the Client.
type Option func(*APIMocker)

// WithDuration returns an Option that sets the sleep duration for the mock.
func WithDuration(duration time.Duration) Option {
	return func(s *APIMocker) {
		s.sleepDuration = &duration
	}
}

// WithDoer returns an Option that sets the Doer for the mock.
func WithDoer(doer func(req *http.Request) (*http.Response, error)) Option {
	return func(s *APIMocker) {
		s.Doer = doer
	}
}

// NewMocker returns a new APIMocker.
func NewMocker(opts ...Option) *APIMocker {
	aMocker := &APIMocker{
		Getter: func(url string) (resp *http.Response, err error) {
			return getMockResponse(), nil
		},
		Poster: func(url string, data url.Values) (resp *http.Response, err error) {
			return getMockResponse(), nil
		},
		Doer: func(req *http.Request) (*http.Response, error) {
			return getMockResponse(), nil
		},
	}
	for _, opt := range opts {
		opt(aMocker)
	}

	return aMocker
}

func getMockResponse() *http.Response {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"ok":"true"}`))),
	}
}
