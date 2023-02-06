package speedtest_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"speed_test_utility/internal/service/speedtest"
	testutils "speed_test_utility/internal/test_utils"
	"testing"
	"time"

	speedtestLib "github.com/abergasov/speedtest-go/speedtest"
	"github.com/stretchr/testify/require"
)

func TestSpeedTest_Execute(t *testing.T) {
	lastMeasure := time.Duration(0)
	for i := 3; i < 9; i++ {
		sleepDuration := time.Duration(i*10) * time.Millisecond
		svs := getService(t, testutils.WithDuration(sleepDuration))
		before := time.Now()
		require.NoError(t, svs.Prepare())
		require.NoError(t, svs.Execute())
		measure := time.Since(before)
		require.True(t, measure > lastMeasure)
		lastMeasure = measure
		require.NotEmptyf(t, svs.GetResult(), "result is empty")
	}
}

func Benchmark_RunSpeedTest(b *testing.B) {
	svs := getService(b)
	if err := svs.Prepare(); err != nil {
		b.Errorf("failed to prepare speedtest service: %v", err)
	}
	for i := 0; i < b.N; i++ {
		require.NoError(b, svs.Execute())
	}
}

func getService(tb testing.TB, opts ...testutils.Option) *speedtest.SpeedTest {
	benchServers, err := json.Marshal(make([]speedtestLib.Server, 10))
	if err != nil {
		tb.Errorf("failed to marshal servers: %v", err)
	}
	benchUserResp := `<?xml version="1.0" encoding="UTF-8"?><settings><client ip="192.168.1.1" lat="10.0000" lon="20.0000" isp="test" /></settings>`
	benchDoer := func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/speedtest-config.php" {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(benchUserResp))),
			}, nil
		}
		return &http.Response{
			StatusCode:    200,
			ContentLength: 10,
			Body:          io.NopCloser(bytes.NewReader(benchServers)),
		}, nil
	}
	tMock := testutils.NewMocker(append(opts, testutils.WithDoer(benchDoer))...)
	return speedtest.CreateSpeedTester(speedtest.WithCustomClient(tMock))
}
