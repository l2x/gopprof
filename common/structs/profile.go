package structs

import "encoding/gob"

func init() {
	gob.Register(ProfileData{})
}

type ProfileData struct {
	NodeID  string
	Created int64
}
