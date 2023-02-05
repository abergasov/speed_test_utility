package service

//go:generate mockgen -source=abstract.go -destination=abstract_mock.go -package=service

// Tester is an interface for testing internet connection speed using differed providers
type Tester interface {
	Prepare() error
	Execute() error
	GetResult() string
}
