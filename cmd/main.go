package main

import (
	"log"
	"time"

	"github.com/ArteShow/Calculator/application"
	"github.com/ArteShow/Calculator/internal"
)

func main() {
	log.Println("🚀 Starting both servers...")

	// Start internal server in parallel
	go func() {
		internal.RunServerAgent()
	}()

	// Wait to ensure internal server starts first
	time.Sleep(2 * time.Second)

	// Start application server
	application.RunServer()
}
