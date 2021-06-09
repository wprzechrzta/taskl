package commands

type ArgRunner interface {
	Init([]string) error
	Run() error
	Name() string
}
