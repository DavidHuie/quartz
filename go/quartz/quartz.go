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
	Registry *Registry
}

func newQuartz() *Quartz {
	return &Quartz{newRegistry()}
}

// Exports a struct via RPC and generates metadata for each of the struct's methods.
func (q *Quartz) RegisterName(name string, s interface{}) error {
	q.Registry.data[name] = newStructMetadata(s)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		// TODO: only export methods with JSON serializable arguments
		// and responses.
		metadata := &methodMetadata{
			method,
			structFieldToType(method.Type.In(1)),
		}
		q.Registry.data[name].NameToMethodMetadata[method.Name] = metadata
	}
	return rpc.RegisterName(name, s)
}

func (q *Quartz) Start(socketPath string) {
	// Start the server and accept connections on a
	// UNIX domain socket.
	rpc.RegisterName("Quartz", q.Registry)
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
