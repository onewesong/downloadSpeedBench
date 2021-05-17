package main

import "testing"

func TestGetServerList(t *testing.T) {
	hosts := GetServerList()
	t.Log(hosts)
}
