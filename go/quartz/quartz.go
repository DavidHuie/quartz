package quartz

import (
	"encoding/json"
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

type Quartz struct {
	Registry map[string]*StructMetadata
}

type MethodMetadata struct {
	Method         reflect.Method `json:"-"`
	ArgumentToType map[string]string
}

type StructMetadata struct {
	NameToMethodMetadata map[string]*MethodMetadata
	TargetStruct         interface{} `json:"-"`
}

func NewStructMetadata(targetStruct interface{}) *StructMetadata {
	return &StructMetadata{
		make(map[string]*MethodMetadata),
		targetStruct,
	}
}

func RegisterName(name string, s interface{}) error {
	// TODO: check that s is a pointer

	quartz.Registry[name] = NewStructMetadata(s)

	t := reflect.TypeOf(s)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)

		// TODO: only export methods with JSON serializable arguments
		// and responses.

		metadata := &MethodMetadata{
			method,
			StructFieldToType(method.Type.In(1)),
		}

		quartz.Registry[name].NameToMethodMetadata[method.Name] = metadata
	}

	return nil
}

func StructFieldToType(t reflect.Type) map[string]string {
	fieldToType := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		fieldToType[t.Field(i).Name] = t.Field(i).Type.String()
	}

	return fieldToType
}

type CallArgs struct {
	StructName string
	Method     string
	MethodArgs string
}

func (q *Quartz) Call(args *CallArgs, response *interface{}) error {

	// TODO: validate args

	metadata := quartz.Registry[args.StructName]
	method := metadata.NameToMethodMetadata[args.Method].Method

	structValue := reflect.ValueOf(metadata.TargetStruct)
	methodArgsValue := reflect.ValueOf([]byte(args.MethodArgs))
	unmarshallerValue := reflect.ValueOf(json.Unmarshal)
	functionResponse := reflect.Indirect(reflect.ValueOf(response))

	// Determine what arguments the function requires
	argsType := method.Type.In(1)
	responseType := method.Type.In(2).Elem()
	responseValuePointer := reflect.New(responseType)

	// Create a value that's a direct reference to the arg argument
	argStructPointer := reflect.New(argsType)
	argStruct := reflect.Indirect(argStructPointer)

	// Unmarshall the argument json
	// TODO: check for unmarshalling errors
	unmarshallerValue.Call([]reflect.Value{methodArgsValue, argStructPointer})

	// Call the method
	// TODO: check for errors
	method.Func.Call([]reflect.Value{structValue, argStruct, responseValuePointer})

	// Set this method's response value
	rpcResponse := reflect.Indirect(responseValuePointer)
	functionResponse.Set(rpcResponse)

	return nil
}

func (q *Quartz) GetMetadata(_ interface{}, value *map[string]*StructMetadata) error {
	*value = q.Registry
	return nil
}

func init() {
	var err error
	listener, err = net.Listen("unix", os.Getenv("QUARTZ_SOCKET"))
	if err != nil {
		panic(err)
	}

	quartz = &Quartz{}
	quartz.Registry = make(map[string]*StructMetadata)

	rpc.Register(quartz)

	// Cleanup the socket file when the server is killed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-sigc
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
		if err != nil {
			panic(err)
		}
		go jsonrpc.ServeConn(conn)
	}
}
