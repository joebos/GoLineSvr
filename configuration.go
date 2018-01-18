package main

// Settings type for holding all configruation settings
type Settings struct {
	host                 string
	port                 int
	numChildProcessMax   int
	numConnsMax          int
	ControlPort          int
	numLinesPerIndexPage int
}

// NewSetttings - all the configruation settings
func NewSetttings() *Settings {
	setting := new(Settings)
	setting.host = "0.0.0.0"
	setting.port = 10497
	//	setting.numChildProcessMax = 100
	//	setting.numConnsMax = 1024
	//	setting.ControlPort = 8080
	setting.numLinesPerIndexPage = 10000
	return setting
}
