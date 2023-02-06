package fast_test

import (
	"speed_test_utility/internal/service/fast"
	testutils "speed_test_utility/internal/test_utils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testURL = "http://localhost:8080"
)

// check taht on increasing response delay - result speed also increace
func TestFastTest_Execute(t *testing.T) {
	lastMeasure := time.Duration(0)
	lastDownload := float64(0)
	lastUpload := float64(0)
	for i := 3; i < 9; i++ {
		sleepDuration := time.Duration(i*100) * time.Millisecond
		svs := getService(
			testutils.NewMocker(
				testutils.WithDuration(sleepDuration),
			),
		)
		before := time.Now()
		require.NoError(t, svs.Execute())
		measure := time.Now().Sub(before)
		require.True(t, measure > lastMeasure)
		lastMeasure = measure
		require.NotEmptyf(t, svs.GetResult(), "result is empty")
		if lastDownload > 0 && lastUpload > 0 {
			require.True(t, svs.DownloadSpeed < lastDownload)
			require.True(t, svs.UploadSpeed < lastUpload)
		}
		lastDownload = svs.DownloadSpeed
		lastUpload = svs.UploadSpeed
	}
}

func Benchmark_RunFastTest(b *testing.B) {
	svs := getService(testutils.NewMocker())
	for i := 0; i < b.N; i++ {
		require.NoError(b, svs.Execute())
	}
}

func getService(tMock *testutils.APIMocker) *fast.Fast {
	return fast.CreateSpeedTester(fast.WithCustomClient(tMock), fast.WithURL(testURL))
}
