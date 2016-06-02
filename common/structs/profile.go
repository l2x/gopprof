package structs

// ProfileOption is options for profiling.
type ProfileOption struct {
	Name  string
	Sleep int
	Debug int
	GC    bool
	Tmp   string
}

// NewProfileOption return ProfileOption with default value.
func NewProfileOption(name string) ProfileOption {
	return ProfileOption{
		Name:  name,
		Sleep: 30,
		Debug: 1,
		Tmp:   "/tmp",
	}
}

// ProfileData is data of profiling
type ProfileData struct {
	NodeID  string
	Type    string
	Files   map[string]string
	Created int64
}
