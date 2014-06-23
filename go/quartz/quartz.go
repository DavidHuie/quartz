package quartz

import (
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
	return rpc.RegisterName(name, s)
}

func StructFieldToType(t reflect.Type) map[string]string {
	fieldToType := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		fieldToType[t.Field(i).Name] = t.Field(i).Type.String()
	}
	return fieldToType
}

func (q *Quartz) GetMetadata(_ interface{}, value *map[string]*StructMetadata) error {
	*value = q.Registry
	return nil
}

func init() {
	socket_path := os.Getenv("QUARTZ_SOCKET")
	if socket_path == "" {
		socket_path = "/tmp/quartz.socket"
	}

	var err error
	listener, err = net.Listen("unix", socket_path)
	if err != nil {
		panic(err)
	}

	quartz = &Quartz{}
	quartz.Registry = make(map[string]*StructMetadata)

	rpc.Register(quartz)

	// Cleanup the socket file when the server is killed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGTERM)
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
