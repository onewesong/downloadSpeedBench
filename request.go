package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const requestTemplate = "GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: curl/7.54.0\r\nConnection: close\r\nAccept: */*\r\n\r\n"

// 绑定源IP发起http GET请求, 返回接收的响应大小
// 若源IP为空, 则自动设置
func Request(host, uri string, localIP string, timeout time.Duration) (int, error) {
	dialer := net.Dialer{
		LocalAddr: &net.TCPAddr{IP: net.ParseIP(localIP)},
		Timeout:   timeout,
	}
	conn, err := dialer.Dial("tcp", host)
	if err != nil {
		return 0, err
	}
	req := fmt.Sprintf(requestTemplate, uri, host)
	_, err = conn.Write([]byte(req))
	if err != nil {
		return 0, err
	}
	return readFully(conn, timeout)
}

// 接收所有数据
// 同时限制每次读的超时时间
func readFully(conn net.Conn, timeout time.Duration) (int, error) {
	defer conn.Close()
	var buf [512]byte
	var n, size int
	var err error
	ch := make(chan bool)
	read := func() {
		n, err = conn.Read(buf[0:])
		ch <- true
	}
loop:
	go read()
	select {
	case <-ch:
		size += n
		if err != nil {
			if err == io.EOF {
				break
			}
			return size, err
		}
		goto loop
	case <-time.After(timeout):
		return size, errors.New("ReadTimeout")
	}
	return size, nil
}
