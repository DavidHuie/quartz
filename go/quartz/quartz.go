package quartz

// This holds information about exported structs.
type Quartz struct {
	Registry registry
}

func newQuartz() *Quartz {
	return &Quartz{newRegistry()}
}

// Returns the struct registry. This method is exported via RPC
// so that the Ruby client can have knowledge about which structs and
// which methods are exported.
func (q *Quartz) GetMetadata(_ interface{}, value *map[string]*structMetadata) error {
	*value = q.Registry
	return nil
}

var (
	quartz     = newQuartz()
	socketPath = "/tmp/quartz.socket"
)
