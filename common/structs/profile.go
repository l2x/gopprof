package structs

// ProfileData is data of profiling
type ProfileData struct {
	ID      int64
	NodeID  string
	Type    string
	Created int64
	File    string

	// status
	Status int // 0 - pending, 1 - success, 2 - failed
	ErrMsg string

	// option
	Sleep int
	Debug int
	GC    bool
}

// NewProfileData .
func NewProfileData() *ProfileData {
	return &ProfileData{}
}
