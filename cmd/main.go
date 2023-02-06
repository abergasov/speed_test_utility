package main

import (
	"flag"
	"log"
	"speed_test_utility/internal/service"
	"speed_test_utility/internal/service/fast"
	"speed_test_utility/internal/service/speedtest"
)

var (
	provider = flag.String("provider", service.FastProvider, "Provider for test speed")
)

func main() {
	flag.Parse()
	if *provider == "" {
		log.Fatal("provider is not set")
	}
	measurer := service.NewInternetConnectionTester(
		speedtest.CreateSpeedTester(),
		fast.CreateSpeedTester(),
	)
	if err := measurer.Run(*provider); err != nil {
		log.Fatal("failed to run test: ", err)
	}
}
