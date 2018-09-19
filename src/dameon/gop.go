package dameon

import (
	"errors"
	"io/ioutil"
	"log"

	"../models"
)

type GOP struct {
	dameon *Dameon
}

func NewGOP(dameon *Dameon) *GOP {
	return &GOP{
		dameon: dameon,
	}
}

func (s *GOP) SetServiceEnabled(args *models.ServiceStatus, reply *models.ServiceStatus) error {
	iconf, ok := s.dameon.pm.ServiceConfigs[args.Name]
	if !ok {
		return errors.New("service does not exist")
	}

	iconf.Enabled = args.Enabled
	iconf.Save(args.Name)

	*reply = *args
	return nil
}

func (s *GOP) RemoveService(serviceName *string, reply *string) error {
	_, ok := s.dameon.pm.ServiceConfigs[*serviceName]
	if !ok {
		return errors.New("service does not exist")
	}
	delete(s.dameon.pm.ServiceConfigs, *serviceName)

	*reply = *serviceName
	return DeleteServiceFile(*serviceName)
}

func (s *GOP) StartService(pargs *models.ProcessStartArgs, reply *models.ProcessStartInfo) error {
	i := pargs.InstanceName
	log.Println(pargs)

	iconf, ok := s.dameon.pm.ServiceConfigs[i]
	if !ok {
		var err error

		iconf, err = s.dameon.pm.NewOrExistingService(i, pargs.Path, pargs.Args)
		if err != nil {
			return err
		}
	}

	p, err := s.dameon.pm.LaunchProcess(i, iconf.Path, iconf.Args)
	pid := 0

	if p.cmd.Process != nil {
		pid = p.cmd.Process.Pid
	}

	r := models.ProcessStartInfo{
		Success:      err == nil,
		PID:          pid,
		InstanceName: p.InstanceName,
	}
	*reply = r
	return err
}

func (s *GOP) StartProcess(pargs *models.ProcessStartArgs, reply *models.ProcessStartInfo) error {
	p, err := s.dameon.pm.LaunchProcess(pargs.InstanceName, pargs.Path, pargs.Args)
	pid := 0

	if p.cmd.Process != nil {
		pid = p.cmd.Process.Pid
	}

	r := models.ProcessStartInfo{
		Success:      err == nil,
		PID:          pid,
		InstanceName: p.InstanceName,
	}

	*reply = r
	return err
}

func (s *GOP) StopProcess(instanceName *string, _ *int) error {
	return s.dameon.pm.StopInstance(*instanceName)
}

func (s *GOP) List(_ *int, reply *models.CurrentListInfo) error {
	clist := models.CurrentListInfo{}

	clist.Processes = make(map[string]string)
	for instance, p := range s.dameon.pm.Processes {
		clist.Processes[instance] = p.cmd.Path
	}

	clist.Services = make(map[string]bool)
	for instance, p := range s.dameon.pm.ServiceConfigs {
		clist.Services[instance] = p.Enabled
	}

	*reply = clist
	return nil
}

func (s *GOP) ReadLog(_ int, log *models.DameonLogInfo) error {
	dl := models.DameonLogInfo{
		FileName: "none",
	}

	if !s.dameon.loggingToFile {
		*log = dl
		return nil
	}

	d, err := ioutil.ReadFile(s.dameon.logFileName)

	dl.FileName = s.dameon.logFileName
	dl.Content = string(d)
	*log = dl
	if err != nil {
		panic(err)
	}
	return err
}
