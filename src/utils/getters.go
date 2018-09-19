package utils

import "fmt"

func GetClientInterface() string {
	return fmt.Sprintf("%s:%d", Options.TargetDameon, Options.InterfacePort)
}

func GetDameonInterface() string {
	bindInterface := ""
	if Options.ExposeDameonInterface {
		bindInterface = "0.0.0.0"
	}
	return fmt.Sprintf("%s:%d", bindInterface, Options.InterfacePort)
}
