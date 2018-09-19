package utils

import (
	"github.com/jessevdk/go-flags"
)

// Options for gop
var Options struct {
	Name string `short:"n" description:"name of service/process instance"`

	RelitivePath bool `short:"r" long:"relitive-path" description:"Will not make path absolute."`

	Service bool `short:"s" description:"If the command is to control a service or an instance."`

	// Append log file
	DisableLog int `long:"disable-log" description:"Do not log dameon."`

	// Dameon rpc listening port
	InterfacePort int `short:"p" long:"port" description:"Port for rpc server to listen on in dameon mode, or client connection in client mode." default:"58823"`

	// Target dameon
	TargetDameon string `short:"t" long:"target-dameon" description:"Dameon service interface to connect to. " default:"127.0.0.1"`

	// Path of the file to contail process id if the current dameon
	LockFile string `long:"lock-file" description:"Path of the dameon lock file." default:"dameon.lock"`

	// If to run the program as a dameon or as a service.
	DameonInstance bool `short:"d" long:"dameon-instance" description:"Launches gop in dameon mode."`

	// If to publicly expose rpc interface for access outside the machine.
	ExposeDameonInterface bool `long:"expose-dameon" description:"Bind dameon to a network interface. Can be dangerious."`

	LogDirectory      string `long:"log-dir" description:"where to store log files" default:"logs"`
	ServicesDirectory string `long:"service-dir" description:"where to store services files" default:"services"`

	Verbose bool `short:"v" description:"Verbose output"`
}

// ParseFlags parses the flags...
func ParseFlags() ([]string, error) {
	return flags.Parse(&Options)
}
