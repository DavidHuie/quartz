package main

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

var (
	quartz   *Quartz
	listener net.Listener
)

type Test struct{}

type Args struct {
	A int
	B int
}

func (q *Test) Add(args Args, val *int) error {
	log.Printf("Received add request, args: %s", args)
	*val = args.A + args.B
	return nil
}

type Quartz struct {
	Registry map[string]*StructMetadata
}

type StructMetadata struct {
	S            interface{}               `json:"-"`
	MethodNames  []string                  `json:"method_names"`
	NameToMethod map[string]reflect.Method `json:"-"`
}

func Register(name string, s interface{}) {
	quartz.Registry[name] =
		&StructMetadata{s, make([]string, 0), make(map[string]reflect.Method)}

	t := reflect.TypeOf(s)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		quartz.Registry[name].MethodNames =
			append(quartz.Registry[name].MethodNames, method.Name)
		quartz.Registry[name].NameToMethod[method.Name] = method
	}
}

type CallArgs struct {
	StructName string
	Method     string
	MethodArgs string
}

func (q *Quartz) Call(args *CallArgs, response *interface{}) error {
	metadata := quartz.Registry[args.StructName]
	method := metadata.NameToMethod[args.Method]
	structValue := reflect.ValueOf(metadata.S)
	methodArgsValue := reflect.ValueOf([]byte(args.MethodArgs))
	argsType := method.Type.In(1)
	responseType := method.Type.In(2).Elem()

	// Create empty struct of argsType
	argStructPointer := reflect.New(argsType)
	argStruct := reflect.Indirect(argStructPointer)

	// Create empty struct of responseType
	responseValuePointer := reflect.New(responseType)

	// Deserialize JSON into empty struct
	unmarshaller := reflect.ValueOf(json.Unmarshal)

	// Todo: check for errors here
	unmarshaller.Call([]reflect.Value{methodArgsValue, argStructPointer})

	// Make call
	method.Func.Call([]reflect.Value{structValue, argStruct, responseValuePointer})

	// Add result to function response
	functionResponse := reflect.Indirect(reflect.ValueOf(response))
	rpcResponse := reflect.Indirect(responseValuePointer)
	functionResponse.Set(rpcResponse)

	return nil
}

func (q *Quartz) GetMetadata(_ interface{}, value *string) error {
	bytes, err := json.Marshal(q.Registry)
	if err != nil {
		return err
	}
	*value = string(bytes)
	return nil
}

func init() {
	var err error
	listener, err = net.Listen("unix", "/tmp/quartz.sock")
	if err != nil {
		panic(err)
	}

	quartz = &Quartz{}
	quartz.Registry = make(map[string]*StructMetadata)

	rpc.Register(quartz)

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
}

func Start() {
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

func main() {
	t := &Test{}
	Register("adder", t)
	Start()
}
