package quartz

type Registry struct {
	data map[string]*structMetadata
}

func newRegistry() *Registry {
	r := &Registry{}
	r.data = make(map[string]*structMetadata)
	return r
}

// Returns the entire registry. This method is exported via RPC
// so that the Ruby client can have knowledge about which structs and
// which methods are exported.
func (r *Registry) GetMetadata(_ interface{}, value *map[string]*structMetadata) error {
	*value = r.data
	return nil
}
