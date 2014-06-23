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

// This holds information about exported structs.
type Quartz struct {
	Registry map[string]*StructMetadata
}

// Return's the struct registry. This method is exported via RPC
// so that the Ruby client can have knowledge about which structs and
// which methods are exported.
func (q *Quartz) GetMetadata(_ interface{}, value *map[string]*StructMetadata) error {
	*value = q.Registry
	return nil
}

var (
	quartz = &Quartz{
		Registry: make(map[string]*StructMetadata),
	}
	socketPath = "/tmp/quartz.socket"
)

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

type MethodMetadata struct {
	Method         reflect.Method `json:"-"`
	ArgumentToType map[string]string
}

// Exports a struct via RPC and generates metadata for each of the struct's methods.
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

// Given a struct type, creates a mapping of field name
// to string representation of the field name's type.
func StructFieldToType(t reflect.Type) map[string]string {
	fieldToType := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		fieldToType[t.Field(i).Name] = t.Field(i).Type.String()
	}
	return fieldToType
}

func Start() {
	// Start the server and accept connections on a
	// UNIX domain socket.
	rpc.Register(quartz)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go jsonrpc.ServeConn(conn)
	}

	// Destroy the socket file when the server is killed.
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

func init() {
	// The Ruby gem sets this environment variable for us.
	if os.Getenv("QUARTZ_SOCKET") != "" {
		socketPath = os.Getenv("QUARTZ_SOCKET")
	}
}
