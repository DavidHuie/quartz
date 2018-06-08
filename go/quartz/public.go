package quartz

var (
	quartz = newQuartz()
)

// RegisterName registers a struct with Quartz by name. The public
// methods on this struct are what become available to Ruby.
func RegisterName(name string, s interface{}) error {
	return quartz.RegisterName(name, s)
}

// Start starts the local Quartz server. This should be the last
// function called by an executable that uses Quartz.
func Start() {
	quartz.Start()
}
