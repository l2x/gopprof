package client

type Option struct {
	Interval int
}

func NewOption() *Option {
	return &Option{}
}
