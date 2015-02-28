package main

import "github.com/DavidHuie/quartz/go/quartz"

type Example struct{}

type ExampleArgs struct {
	Num int
}

type Response struct {
	Output []string
}

func (t *Example) LongJson(args ExampleArgs, response *Response) error {
	result := []string{}
	for i := 0; i < 8192; i++ {
		result = append(result, "Test")
	}
	*response = Response{result}
	return nil
}

func main() {
	example := &Example{}
	quartz.RegisterName("example", example)
	quartz.Start()
}
