package main

import (
	"flag"
	"fmt"
	"golang.org/x/sync/semaphore"
)

type PortScanner struct {
	ip string
	lock *semaphore.Weighted
}

func getCMDData() (int64, string) {
	limit := flag.Int64("n", 10, "Numbe of concurrent threads")
	ip := flag.String("ip", "", "IP address in the network to scan")
	flag.Parse()
	return *limit, *ip
}

func main() {
	limit, ip := getCMDData()
	
	ps := &PortScanner{
		ip: ip,
		lock: semaphore.NewWeighted(limit),
	}

	fmt.Printf("IP = %s and Limit = %d", ps.ip, limit)
}