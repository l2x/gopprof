package structs

import "encoding/gob"

func init() {
	gob.Register(ProfileData{})
}

type ProfileData struct {
	NodeID  string
	Created int64

	Type      string
	File      string
	BinMD5    string
	GoVersion string
}

type Goroot struct {
	Version string `json:"version"`
	Path    string `json:"path"`
}
