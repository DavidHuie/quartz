package quartz

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// Quartz holds information about exported structs.
type Quartz struct {
	registry *registry
}

func newQuartz() *Quartz {
	return &Quartz{newRegistry()}
}

// RegisterName exports a struct via RPC and generates metadata for
// each of the struct's methods.
func (q *Quartz) RegisterName(name string, s interface{}) error {
	meta := newStructMetadata(s)
	t := reflect.TypeOf(s)

	// Register all public methods.
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		metadata := &methodMetadata{
			method,
			structFieldToType(method.Type.In(1)),
		}
		meta.NameToMethodMetadata[method.Name] = metadata
	}

	q.registry.addMetadata(name, meta)

	return rpc.RegisterName(name, s)
}

// Start launches the server.
func (q *Quartz) Start() {
	log.SetOutput(os.Stderr)

	// The Ruby gem sets this environment variable for us.
	socketPath := "/tmp/quartz.socket"
	if os.Getenv("QUARTZ_SOCKET") != "" {
		socketPath = os.Getenv("QUARTZ_SOCKET")
	}

	// Start the server and accept connections on a
	// UNIX domain socket.
	rpc.RegisterName("Quartz", q.registry)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}

	// Close the socket when the server is killed.
	quit := make(chan struct{})
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigc
		if err := listener.Close(); err != nil {
			log.Printf("error closing listener: %s", err)
		}

		close(quit)
	}()

Loop:
	for {
		select {
		case <-quit:
			break Loop
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %s", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}

	os.Exit(0)
}
