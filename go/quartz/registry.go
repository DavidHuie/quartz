package quartz

import "sync"

// registry contains metadata about all rpc-exported structs.
type registry struct {
	mutex *sync.Mutex
	data  map[string]*structMetadata
}

func (r *registry) addMetadata(name string, s *structMetadata) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.data[name] = s
}

func newRegistry() *registry {
	r := &registry{
		mutex: &sync.Mutex{},
	}
	r.data = make(map[string]*structMetadata)
	return r
}

// GetMetadata Returns the entire registry. This method is exported
// via RPC so that the Ruby client can have knowledge about exported
// structs and methods.
func (r *registry) GetMetadata(_ interface{}, value *map[string]*structMetadata) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	*value = r.data
	return nil
}
