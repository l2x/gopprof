package util

import "testing"

func TestInStringSlice(t *testing.T) {
	s := "a"
	arr := []string{"1", "a", "b"}
	if InStringSlice(s, arr) == false {
		t.Error("string should be in slice")
	}
	s = "c"
	if InStringSlice(s, arr) == true {
		t.Error("string should not be in slice")
	}
}
