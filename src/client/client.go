package client

import (
	"log"
	"net/rpc"

	"../models"
	"../utils"
)

type Client struct {
	rpc *rpc.Client
}

func New() (*Client, error) {
	var err error
	c := new(Client)

	c.rpc, err = rpc.DialHTTP("tcp", utils.GetClientInterface())

	return c, err
}

//Process commands for client
func (c *Client) ProcessCommands(args []string) {
	if len(args) == 1 { //Commands with no arguments
		switch args[0] {
		case "log":
			c.ReadLog()
		case "list":
			c.List()
		case "start":
			c.Start("", nil)
		}
	}

	if len(args) > 1 {
		switch args[0] {
		case "start":
			c.Start(args[1], args[2:])
		case "stop":
			c.Stop(args[1])
		case "enable":
			c.SetServiceEnabled(args[1], true)
		case "disable":
			c.SetServiceEnabled(args[1], false)
		case "remove":
			c.RemoveService(args[1])
		}
	}
}

func (c *Client) ReadLog() {
	logInfo := new(models.DameonLogInfo)

	if err := c.rpc.Call(gopReadLog, 0, logInfo); err != nil {
		utils.LogError(err)
	}
	log.Println("[", logInfo.FileName, "]\n\r", logInfo.Content)
}
