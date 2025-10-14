package app

type Application interface {
	Run() error
	Shutdown() error
}
