package quartz

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
)

// Quartz holds information about exported structs.
type Quartz struct {
	Registry *Registry
}

func newQuartz() *Quartz {
	return &Quartz{newRegistry()}
}

// RegisterName exports a struct via RPC and generates metadata for
// each of the struct's methods.
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

// Start launches the server.
func (q *Quartz) Start() {
	// The Ruby gem sets this environment variable for us.
	socketPath := "/tmp/quartz.socket"
	if os.Getenv("QUARTZ_SOCKET") != "" {
		socketPath = os.Getenv("QUARTZ_SOCKET")
	}

	// Start the server and accept connections on a
	// UNIX domain socket.
	rpc.RegisterName("Quartz", q.Registry)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}

	// Destroy the socket file when the server is killed.
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		<-sigc
		listener.Close()
		os.Remove(socketPath)
		os.Exit(0)
	}()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			// The connection was closed by a signal.
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			fmt.Fprintf(os.Stderr, "error accepting connection: %s", err)
			continue
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Fprintf(os.Stderr, "Recovered from panic: %s", r)
				}
			}()

			jsonrpc.ServeConn(conn)
		}()
	}
}
