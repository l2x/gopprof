package client

type ClientOption struct {
	Interval int
}

func NewClientOption() *ClientOption {
	return &ClientOption{}
}

type ProfileOption struct {
	Name  string
	Sleep int
	Debug int
	GC    bool
	tmp   string
}

func NewProfileOption(name string) *ProfileOption {
	return &ProfileOption{
		Name:  name,
		Sleep: 30,
		Debug: 1,
		tmp:   "/tmp",
	}
}
