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

// Exports a struct via RPC and generates metadata for each of the struct's methods.
func RegisterName(name string, s interface{}) error {
	quartz.Registry[name] = newStructMetadata(s)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		// TODO: only export methods with JSON serializable arguments
		// and responses.
		metadata := &methodMetadata{
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
	// The Ruby gem sets this environment variable for us.
	if os.Getenv("QUARTZ_SOCKET") != "" {
		socketPath = os.Getenv("QUARTZ_SOCKET")
	}

	// Start the server and accept connections on a
	// UNIX domain socket.
	rpc.Register(quartz)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
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

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go jsonrpc.ServeConn(conn)
	}
}
