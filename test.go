package main

import (
	"github.com/DavidHuie/quartz/go/quartz"
)

type Test struct{}

type Args struct {
	A int
	B int
}

type Response struct {
	X int
}

func (t *Test) Add(args Args, response *Response) error {
	*response = Response{args.A + args.B}
	return nil
}

func main() {
	a := &Test{}
	quartz.RegisterName("adder", a)
	quartz.Start()
}
