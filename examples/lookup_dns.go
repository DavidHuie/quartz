package main

import (
	"net"

	"github.com/DavidHuie/quartz/go/quartz"
)

type Resolver struct{}

type FindIPsArgs struct {
	Hostnames []string
}

type FindIPsResponse struct {
	HostnameToIPs map[string][]net.IP
}

func (r *Resolver) FindIPs(args FindIPsArgs, response *FindIPsResponse) error {
	*response = FindIPsResponse{}
	response.HostnameToIPs = make(map[string][]net.IP)
	c := make(chan bool)

	for _, hostname := range args.Hostnames {
		go func(h string) {
			addrs, err := net.LookupIP(h)
			if err != nil {
				panic(err)
			}

			response.HostnameToIPs[h] = addrs

			c <- true
		}(hostname)
	}

	for i := 0; i < len(args.Hostnames); i++ {
		<-c
	}

	return nil
}

func main() {
	resolver := &Resolver{}
	quartz.RegisterName("resolver", resolver)
	quartz.Start()
}
