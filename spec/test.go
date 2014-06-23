package main

import (
	"errors"

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

func (t *Test) AddError(args Args, response *Response) error {
	return errors.New("Adding error")
}

func main() {
	a := &Test{}
	quartz.RegisterName("adder", a)
	quartz.Start()
}
