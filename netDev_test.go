package main

import "testing"

func TestMatchDevice(t *testing.T) {
	prefix := "en"
	IPlist, err := MatchDevice(prefix)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(IPlist)
}
