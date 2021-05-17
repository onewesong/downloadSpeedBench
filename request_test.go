package main

import (
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	host := "tn1.chtm.hinet.net:8080"
	uri := "/speedtest/random350x350.jpg"
	localIP := "172.20.120.28"
	timeout := time.Second * 5
	t.Log("start request")
	size, err := Request(host, uri, localIP, timeout)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("read size: ", size)
	if size < 24000 {
		t.Fatal("size too small")
	}
}

func TestRequestWithEmptyLocalIP(t *testing.T) {
	host := "tn1.chtm.hinet.net:8080"
	uri := "/speedtest/random350x350.jpg"
	localIP := ""
	timeout := time.Second * 5
	t.Log("start request")
	size, err := Request(host, uri, localIP, timeout)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("read size: ", size)
	if size < 24000 {
		t.Fatal("size too small")
	}
}
