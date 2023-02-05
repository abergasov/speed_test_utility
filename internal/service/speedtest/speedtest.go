package speedtest

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/abergasov/speedtest-go/speedtest"
)

// SpeedTest is a wrapper for speedtest-go/speedtest
// uses github.com/showwin/speedtest-go/speedtest utility to execute tests and fetch results
// for tests and benchmarks used interfaced http.Client
type SpeedTest struct {
	httpClient speedtest.HTTPClient
	client     *speedtest.Speedtest
	user       *speedtest.User
	serverList speedtest.Servers
}

// Option is a function that can be passed to New to modify the Client.
type Option func(*SpeedTest)

// WithCustomClient sets the http.Client used to make requests.
func WithCustomClient(client speedtest.HTTPClient) Option {
	return func(s *SpeedTest) {
		s.httpClient = client
	}
}

// CreateSpeedTester create wrapper for run speedtest
// uses github.com/showwin/speedtest-go/speedtest utility to execute tests and fetch results
func CreateSpeedTester(opts ...Option) *SpeedTest {
	res := &SpeedTest{
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(res)
	}
	res.client = speedtest.New(speedtest.WithDoer(res.httpClient))
	return res
}

// Prepare fetches user info and server list
func (s *SpeedTest) Prepare() (err error) {
	s.user, err = s.client.FetchUserInfo()
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}
	serverList, err := s.client.FetchServers(s.user)
	if err != nil {
		return fmt.Errorf("failed to fetch servers list: %w", err)
	}
	s.serverList, err = serverList.FindServer([]int{})
	if err != nil {
		return fmt.Errorf("failed to find server: %w", err)
	}
	return nil
}

// Execute executes speedtest for all servers in server list
func (s *SpeedTest) Execute() (err error) {
	for _, srv := range s.serverList {
		log.Printf("executing speedtest, server: %s, host: %s\n", srv.Name, srv.Host)
		log.Println("pinging server")
		if err = srv.PingTest(); err != nil {
			return fmt.Errorf("failed to ping server: %w", err)
		}
		log.Println("executing download test")
		if err = srv.DownloadTest(false); err != nil {
			return fmt.Errorf("failed to execute download test: %w", err)
		}
		log.Println("executing upload test")
		if err = srv.UploadTest(false); err != nil {
			return fmt.Errorf("failed to execute upload test: %w", err)
		}
	}
	return err
}

// GetResult returns string with results of speedtest
func (s *SpeedTest) GetResult() string {
	sb := &strings.Builder{}
	for _, srv := range s.serverList {
		sb.WriteString(
			fmt.Sprintf("server: %s, latency: %s, download: %f, upload: %f", srv.Name, srv.Latency.String(), srv.DLSpeed, srv.ULSpeed),
		)
	}
	return sb.String()
}
