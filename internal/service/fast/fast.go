package fast

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"speed_test_utility/internal/external"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gesquive/fast-cli/fast"
	"golang.org/x/sync/errgroup"
)

const (
	workload      = 8
	payloadSizeMB = 25.0 // download payload is by default 25MB, make upload 25MB also
)

type Fast struct {
	client          external.HTTPClient
	url             string // url to test
	DownloadSpeed   float64
	UploadSpeed     float64
	totalDownloaded int64
	uploadPayload   string // pre-generated upload of payloadSizeMB
}

// Option is a function that can be passed to New to modify the Client.
type Option func(*Fast)

// WithCustomClient sets the custom client to make requests.
func WithCustomClient(client external.HTTPClient) Option {
	return func(s *Fast) {
		s.client = client
	}
}

// WithURL sets the url to test
func WithURL(urls string) Option {
	return func(s *Fast) {
		s.url = urls
	}
}

// CreateSpeedTester create wrapper for run fast
// uses https://github.com/adhocore/fast utility to execute tests and fetch results
func CreateSpeedTester(opts ...Option) *Fast {
	ft := &Fast{
		client: http.DefaultClient,
		// generate 10b * x MB / 10 = x MB payload
		uploadPayload: strings.Repeat("0123456789", payloadSizeMB*1024*1024/10),
	}
	for _, opt := range opts {
		opt(ft)
	}
	return ft
}

// Prepare fetches user info and server list. Uses chrome headless browser to open page and get token
func (s *Fast) Prepare() (err error) {
	atomic.StoreInt64(&s.totalDownloaded, 0)
	urls := fast.GetDlUrls(4)
	if len(urls) == 0 {
		return fmt.Errorf("failed to fetch download urls")
	}
	s.url = urls[0]
	return nil
}

// Execute sequentially executes download and upload tests
func (s *Fast) Execute() (err error) {
	log.Printf("executing fast test, server: %s", s.url)
	log.Println("executing download test")
	s.DownloadSpeed, err = s.measureDownloadSpeed(s.url)
	if err != nil {
		return fmt.Errorf("failed to execute download test: %w", err)
	}
	log.Println("executing upload test")
	s.UploadSpeed, err = s.measureUploadSpeed(s.url)
	if err != nil {
		return fmt.Errorf("failed to execute upload test: %w", err)
	}
	return nil
}

func (s *Fast) GetResult() string {
	return fmt.Sprintf("download: %f, upload: %f", s.DownloadSpeed, s.UploadSpeed)
}

// measureUploadSpeed run operation in parallel and measure speed
func (s *Fast) measureUploadSpeed(urlAddr string) (float64, error) {
	eg := errgroup.Group{}

	sTime := time.Now()
	for i := 0; i < workload; i++ {
		eg.Go(func() error {
			return s.uploadRequest(urlAddr)
		})
	}
	if err := eg.Wait(); err != nil {
		return 0, fmt.Errorf("failed to execute operation: %w", err)
	}

	return payloadSizeMB * 8 * float64(workload) / time.Now().Sub(sTime).Seconds(), nil
}

// measureDownloadSpeed execute download operation in parallel and measure speed
// every request consume all payload and count bytes consumed
func (s *Fast) measureDownloadSpeed(urlAddr string) (float64, error) {
	eg := errgroup.Group{}
	sTime := time.Now()
	for i := 0; i < workload; i++ {
		eg.Go(func() error {
			return s.downloadRequest(urlAddr)
		})
	}
	if err := eg.Wait(); err != nil {
		return 0, fmt.Errorf("failed to measure download speed: %w", err)
	}
	secPassed := time.Now().Sub(sTime).Seconds()
	contentConsumed := atomic.LoadInt64(&s.totalDownloaded)

	reqMB := contentConsumed / 1024 / 1024
	return float64(reqMB) * 8 / secPassed, nil
}

func (s *Fast) downloadRequest(urlAddr string) error {
	resp, err := s.client.Get(urlAddr)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	atomic.AddInt64(&s.totalDownloaded, int64(len(bodyBytes)))
	return nil
}

func (s *Fast) uploadRequest(uri string) error {
	v := url.Values{}
	v.Add("content", s.uploadPayload)

	resp, err := s.client.PostForm(uri, v)
	if err != nil {
		return fmt.Errorf("failed to post form: %w", err)
	}
	defer resp.Body.Close()
	return nil
}
