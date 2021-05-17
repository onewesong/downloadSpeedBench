package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/panjf2000/ants"
)

const (
	bestN = 10
)

var (
	IPListPtr    *[]string
	SIZES        = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000, 4000, 4000, 4000}
	Threads      = flag.Int("w", 20, "并发数")
	Timeout      = flag.Int("t", 3, "单个请求超时时间")
	NetDevPrefix = flag.String("p", "", "网卡前缀")
)

func getBestServer(hosts []string) []string {
	var wg sync.WaitGroup
	wg.Add(bestN)
	ch := make(chan string, bestN)
	for _, host := range hosts {
		go func(host string) {
			url := fmt.Sprintf("http://%s/speedtest/random750x750.jpg", host)
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == 200 {
				ch <- host
				resp.Body.Close()
				wg.Done()
			}
		}(host)
	}
	wg.Wait()
	wg.Add(len(hosts))
	res := []string{}
	for i := 0; i < bestN; i++ {
		res = append(res, <-ch)
	}
	return res
}

func download(localIP, host string, size int, ch chan bool) {
	ch <- true
	fmt.Print("S")
	uri := fmt.Sprintf("/speedtest/random%vx%v.jpg", size, size)
	timeout := time.Duration(*Timeout) * time.Second
	_, err := Request(host, uri, localIP, timeout)
	if err != nil {
		fmt.Print("F")
	}
	fmt.Print(".")
	<-ch
}

func main() {
	flag.Parse()
	if *NetDevPrefix != "" {
		IPList, err := MatchDevice(*NetDevPrefix)
		if err != nil {
			log.Fatal(err)
		}
		if len(IPList) == 0 {
			log.Fatal("获取不到匹配的网卡前缀, 请确认是否有权限")
		}
		log.Println("匹配到的网卡IP: ", IPList)
		IPListPtr = &IPList
	} else {
		IPListPtr = &[]string{""}
	}
	log.Printf("参数: Threads: %v, Timeout: %v", *Threads, *Timeout)
	serverHosts := GetServerList()
	log.Printf("got %v servers", len(serverHosts))
	hosts := getBestServer(serverHosts)
	log.Println("select best", len(hosts), "servers: ", hosts)
	log.Println("start test download speed")
	ch := make(chan bool, *Threads)
	defer close(ch)

	// 创建goroutine池, 限制创建的goroutine池在一定数量且循环复用
	p, _ := ants.NewPool(1000)
	defer p.Release()
	for {
		fmt.Print("R")
		for _, ip := range *IPListPtr {
			for _, host := range hosts {
				for _, size := range SIZES {
					_ = p.Submit(func() {
						download(ip, host, size, ch)
					})
				}
			}
		}
	}
}
