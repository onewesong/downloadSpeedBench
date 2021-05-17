package main

import (
	"strings"

	"github.com/google/gopacket/pcap"
)

// 返回匹配指定网卡前缀的IPv4地址列表
func MatchDevice(devNamePrefix string) ([]string, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	IPList := []string{}
	for _, device := range devices {
		if strings.HasPrefix(device.Name, devNamePrefix) {
			for _, address := range device.Addresses {
				ipAddr := address.IP.String()
				// 排除IPv6
				if !strings.Contains(ipAddr, ":") {
					IPList = append(IPList, ipAddr)
				}
			}
		}
	}
	return IPList, nil
}
