package utils

import (
	"log"
	"os"
)

func LogVoberseErrorFatal(err error) {
	if Options.Verbose {
		log.Fatal(err)
	} else {
		os.Exit(1)
	}
}

func LogError(err error) {
	if Options.Verbose {
		log.Panic(err)
	} else {
		log.Fatal(err)
	}
}
