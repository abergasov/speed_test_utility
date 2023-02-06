package service

import (
	"fmt"
	"log"
)

const (
	SpeedTestProvider = "speedtest"
	FastProvider      = "fast"
)

// InternetConnectionTester orchestrates internet connection tests
type InternetConnectionTester struct {
	mapper map[string]Tester
}

// NewInternetConnectionTester creates new instance of InternetConnectionTester which orchestrates tests
func NewInternetConnectionTester(speedTester, fastTester Tester) *InternetConnectionTester {
	return &InternetConnectionTester{
		mapper: map[string]Tester{
			SpeedTestProvider: speedTester,
			FastProvider:      fastTester,
		},
	}
}

// Run executes test for given provider
func (i *InternetConnectionTester) Run(provider string) error {
	if _, ok := i.mapper[provider]; !ok {
		return fmt.Errorf("unknown provider: %s", provider)
	}

	log.Println("starting test")
	if err := i.mapper[provider].Prepare(); err != nil {
		return fmt.Errorf("failed to prepare test: %w", err)
	}
	log.Println("test prepared, execute it")
	if err := i.mapper[provider].Execute(); err != nil {
		return fmt.Errorf("failed to execute test: %w", err)
	}
	log.Printf("test finished, result: %s\n", i.mapper[provider].GetResult())
	return nil
}
