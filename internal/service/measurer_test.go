package service_test

import (
	"speed_test_utility/internal/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestInternetConnectionTester_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	fastMock := service.NewMockTester(ctrl)
	speedtestMock := service.NewMockTester(ctrl)

	measurer := service.NewInternetConnectionTester(speedtestMock, fastMock)

	t.Run("speedtest", func(t *testing.T) {
		speedtestMock.EXPECT().Prepare().Return(nil).Times(1)
		speedtestMock.EXPECT().Execute().Return(nil).Times(1)
		speedtestMock.EXPECT().GetResult().Return("abc").Times(1)
		require.NoError(t, measurer.Run(service.SpeedTestProvider))
	})

	t.Run("fast", func(t *testing.T) {
		fastMock.EXPECT().Prepare().Return(nil).Times(1)
		fastMock.EXPECT().Execute().Return(nil).Times(1)
		fastMock.EXPECT().GetResult().Return("abc").Times(1)
		require.NoError(t, measurer.Run(service.FastProvider))
	})

	t.Run("unknown provider", func(t *testing.T) {
		require.Error(t, measurer.Run("unknown"))
	})
}
