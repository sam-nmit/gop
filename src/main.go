package main

import (
	"log"

	"./utils"

	"./client"
	"./dameon"

	"github.com/common-nighthawk/go-figure"
)

var (
	gopTitle = figure.NewFigure("GOP", "isometric1", true)

	d dameon.Dameon
)

func main() {
	args, err := utils.ParseFlags()

	if err != nil {
		log.Println("Failed to parse command line arguments.")
		utils.LogError(err)
	}

	if utils.Options.DameonInstance {
		gopTitle.Print()

		log.Println("Starting gop in dameon mode.")

		// Set logging to file.
		if err := d.SetLogging(); err != nil {
			log.Println("Failed to open file for logging.")
			utils.LogVoberseErrorFatal(err)
		}

		//Launch in dameon mode.
		err := d.Start(utils.GetDameonInterface())
		utils.LogError(err)

	} else {
		// Disable timestamp for logging in client mode.
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

		c, err := client.New()
		if err != nil {
			log.Println("Failed to connect to GOP service.")
			utils.LogVoberseErrorFatal(err)
		}

		c.ProcessCommands(args)

		if err != nil {
			utils.LogError(err)
		}
	}

}
