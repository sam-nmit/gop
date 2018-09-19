package client

import (
	"errors"
	"log"
	"path/filepath"

	"../models"
	"../utils"
)

func (c *Client) Start(proc string, args []string) {

	if utils.Options.Service {
		c.StartInstance(gopStartService, proc, args)
	} else {
		if proc == "" {
			utils.LogError(errors.New("must define a process to run"))
		}
		c.StartInstance(gopStartProcess, proc, args)
	}
}
func (c *Client) StartInstance(servicename, proc string, args []string) {
	var err error

	procPath := proc

	if proc != "" && !utils.Options.RelitivePath {
		procPath, err = filepath.Abs(proc)
		if err != nil {
			utils.LogError(err)
		}
	}

	startArgs := &models.ProcessStartArgs{
		Path:         procPath,
		Args:         args,
		InstanceName: utils.Options.Name,
	}

	result := new(models.ProcessStartInfo)
	err = c.rpc.Call(servicename, startArgs, result)

	if err != nil {
		log.Println("Failed to start process.")

		utils.LogError(err)
	}

	if result.Success {
		log.Println("Successfully started process. [", result.InstanceName, "]")
		log.Println("PID:", result.PID)
	} else {
		log.Println(result.Error)
	}
}

func (c *Client) Stop(instance string) {
	var i int
	if err := c.rpc.Call(gopStopProcess, instance, &i); err != nil {
		utils.LogError(err)
	} else {
		log.Println("Instance stopped.")
	}
}

func (c *Client) List() {
	result := new(models.CurrentListInfo)
	err := c.rpc.Call(gopList, 0, result)
	if err != nil {
		utils.LogError(err)
	}
	log.Println("[Processes]")
	for k, p := range result.Processes {
		log.Println("\t", k, " -> ", p)
	}

	log.Println("[Services]")
	for k, p := range result.Services {
		log.Println("\t", k, " -> Enabled: ", p)
	}
}

func (c *Client) SetServiceEnabled(name string, enabled bool) {
	args := &models.ServiceStatus{
		Name:    name,
		Enabled: enabled,
	}

	reply := new(models.ServiceStatus)
	err := c.rpc.Call(gopSetServiceEnabled, args, reply)
	if err != nil {
		utils.LogError(err)
	}

	textStatus := ""
	if reply.Enabled {
		textStatus = "enabled"
	} else {
		textStatus = "disabled"
	}

	log.Println(reply.Name, "is now", textStatus)
}

func (c *Client) RemoveService(name string) {

	var r string
	err := c.rpc.Call(gopRemoveService, &name, &r)
	if err != nil {
		utils.LogError(err)
	}

	log.Println("removed service", r)
}
