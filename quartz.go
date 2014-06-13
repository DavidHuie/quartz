package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"syscall"
)

type Quartz struct{}

type AddArgs struct {
	N int
	M int
}

func (q *Quartz) Add(args AddArgs, val *int) error {
	log.Printf("Received log request, args: %s", args)
	*val = args.N + args.M
	return nil
}

func main() {
	q := &Quartz{}
	rpc.Register(q)

	listener, err := net.Listen("unix", "/tmp/quartz.sock")
	if err != nil {
		panic(err)
	}

	// Clean up the socket file when the server is killed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sigc
		log.Print("Cleaning socket.")
		err := listener.Close()
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		log.Print("Accepted connection.")
		if err != nil {
			panic(err)
		}
		log.Print("Serving connection.")
		jsonrpc.ServeConn(conn)
	}
}
