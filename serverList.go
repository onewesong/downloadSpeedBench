package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type ServerStatic struct {
	XMLName xml.Name       `xml:"settings"`
	Servers Recurlyservers `xml:"servers"`
}

type Recurlyservers struct {
	XMLName xml.Name `xml:"servers"`
	Svs     []server `xml:"server"`
}

type server struct {
	XMLName xml.Name `xml:"server"`
	Host    string   `xml:"host,attr"`
}

func GetServerList() []string {
	resp, err := http.Get("https://www.speedtest.net/speedtest-servers-static.php")
	if err != nil {
		log.Printf("get server list failed: %v. retrying another...", err)
		resp, err = http.Get("https://c.speedtest.net/speedtest-servers-static.php")
		if err != nil {
			log.Fatal("get server list failed: ", err)
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read server list body failed: %s", err)
	}
	v := ServerStatic{}
	err = xml.Unmarshal(body, &v)
	if err != nil {
		log.Fatalf("unmarshal server list body failed: %s", err)
	}
	hosts := []string{}
	for _, s := range v.Servers.Svs {
		hosts = append(hosts, s.Host)
	}
	return hosts
}
