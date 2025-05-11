package main

import (
	"log"

	"github.com/ArteShow/Calculator/application"
	"github.com/ArteShow/Calculator/internal"
	"github.com/ArteShow/Calculator/pkg/setup"
)

func main() {
	log.Println("Starting application...") // Replacing fmt.Println with log
	setup.Setup()
	go func() {
		log.Println("Starting TCP listener...") // Replacing fmt.Println with log
		internal.StartTCPListener()
	}()
	log.Println("Starting HTTP server...") // Replacing fmt.Println with log
	application.StartApplicationServer()
}
