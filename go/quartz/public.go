package quartz

var (
	quartz = newQuartz()
)

func RegisterName(name string, s interface{}) error {
	return quartz.RegisterName(name, s)
}

func Start() {
	quartz.Start()
}
