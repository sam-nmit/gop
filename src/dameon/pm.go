package dameon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"../utils"
)

type ServiceConfig struct {
	Path    string   `json:"path"`
	Args    []string `json:"args"`
	Enabled bool     `json:"enabled"`
}

func (sc *ServiceConfig) Save(name string) error {
	f, err := os.Create(filepath.Join(utils.Options.ServicesDirectory, name))
	if err != nil {
		return err
	}
	defer f.Close()

	j := json.NewEncoder(f)
	j.SetIndent("", "\t")
	j.Encode(sc)
	return nil
}

func DeleteServiceFile(name string) error {
	return os.Remove(filepath.Join(utils.Options.ServicesDirectory, name))
}

type PMProcess struct {
	cmd *exec.Cmd

	InstanceName string
	cxt          context.Context
}

type ProcessExitInfo struct {
	Error error
	Proc  *PMProcess
}

// PM is the process manager.
type PM struct {
	ProcessExits   chan ProcessExitInfo
	ServiceConfigs map[string]ServiceConfig

	Processes map[string]*PMProcess
}

func NewPM() *PM {

	pm := &PM{}
	pm.ProcessExits = make(chan ProcessExitInfo, 1)

	pm.Processes = make(map[string]*PMProcess)

	if err := os.MkdirAll(utils.Options.ServicesDirectory, 0777); err != nil {
		utils.LogError(err)
	}
	pm.ServiceConfigs = make(map[string]ServiceConfig)

	//Load services
	filepath.Walk(utils.Options.ServicesDirectory, pm.loadServiceFile)

	return pm
}

func (pm *PM) loadServiceFile(path string, info os.FileInfo, err error) error {

	if info.IsDir() {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		log.Println("Failed to open service config", path)
		return nil
	}
	defer f.Close()

	servicename := filepath.Base(path)

	j := json.NewDecoder(f)

	sconfig := ServiceConfig{}
	if err := j.Decode(&sconfig); err != nil {
		log.Println("Failed to decode service config", path)
		return nil
	}

	if _, ok := pm.ServiceConfigs[servicename]; ok {
		log.Println("Duplicate service", servicename)
		return nil
	}

	log.Println("Loaded service", servicename)
	pm.ServiceConfigs[servicename] = sconfig

	if sconfig.Enabled {
		_, err = pm.LaunchProcess(servicename, sconfig.Path, sconfig.Args)
		if err != nil {
			log.Println("Failed to start service", servicename, " ->", err)
		}
	}
	return nil
}

func (pm *PM) NewOrExistingService(name, proc string, args []string) (ServiceConfig, error) {
	if s, ok := pm.ServiceConfigs[name]; ok {
		return s, nil
	} else {
		if proc == "" {
			return ServiceConfig{}, errors.New("Service does not exist.")
		}
	}
	sconfig := ServiceConfig{
		Path:    proc,
		Args:    args,
		Enabled: false,
	}
	err := sconfig.Save(name)
	if err != nil {
		return sconfig, errors.New(fmt.Sprint("Failed to create service file for", name))
	}

	pm.ServiceConfigs[name] = sconfig
	return sconfig, nil
}

func (pm *PM) LaunchProcess(instanceName, proc string, args []string) (*PMProcess, error) {
	log.Println("Launching \"", proc, "\" - args ", args)

	if instanceName == "" {
		t := time.Now()
		instanceName = fmt.Sprintf("%s-%d%d-%d%d", filepath.Base(proc), t.Month(), t.Day(), t.Second(), t.Nanosecond()%99)
	}

	cxt := context.Background() //possable timeout for future
	p := &PMProcess{
		cmd:          exec.CommandContext(cxt, proc, args...),
		cxt:          cxt,
		InstanceName: instanceName,
	}

	if _, ok := pm.Processes[instanceName]; ok {
		return nil, errors.New(fmt.Sprint("Process instance name clash", instanceName))
	}
	pm.Processes[instanceName] = p

	err := p.cmd.Start()
	if err == nil {
		go func() {
			err := p.cmd.Wait()
			delete(pm.Processes, instanceName)
			pm.ProcessExits <- ProcessExitInfo{
				Error: err,
				Proc:  p,
			}
		}()
	}

	return p, err
}

func (pm *PM) StopInstance(name string) error {
	i, ok := pm.Processes[name]
	if !ok {
		return errors.New("Instance does not exist")
	}
	return i.cmd.Process.Kill()
}
