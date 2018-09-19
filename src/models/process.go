package models

type ProcessStartArgs struct {
	Path         string
	Args         []string
	InstanceName string
}

type ProcessStartInfo struct {
	Success bool
	Error   string
	PID     int

	InstanceName string
}

type ServiceStatus struct {
	Name    string
	Enabled bool
}
