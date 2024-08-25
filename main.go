package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

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

func scanPort(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			scanPort(ip, port, timeout)
		} else {
			fmt.Printf("Port: %d is closed\n", port)
		}
		return
	}

	conn.Close()
	fmt.Printf("%d - open\n", port)
}

func (ps *PortScanner) start(start, end int, timeout time.Duration) {
	var wg sync.WaitGroup
	defer wg.Wait()

	for port := start; port <= end; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			scanPort(ps.ip, port, timeout)
		}(port)
	}
}

func main() {
	limit, ip := getCMDData()
	
	ps := &PortScanner{
		ip: ip,
		lock: semaphore.NewWeighted(limit),
	}

	fmt.Printf("IP = %s and Limit = %d\n", ps.ip, limit)
	ps.start(1, 65535, 500*time.Millisecond)
}