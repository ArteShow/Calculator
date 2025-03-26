package main

import (
	"log"
	"time"

	"github.com/ArteShow/Calculator/application"
	"github.com/ArteShow/Calculator/internal"
	"github.com/ArteShow/Calculator/ui" // Import the UI package
)

func main() {
	log.Println("🚀 Starting both servers...")

	// Start internal server in parallel
	go func() {
		log.Println("🛠️ Starting internal server...")
		internal.RunServerAgent()
		log.Println("✅ Internal server started!")
	}()

	// Wait to ensure internal server starts first
	time.Sleep(2 * time.Second)

	// Start application server
	go func() {
		log.Println("🛠️ Starting application server...")
		application.RunServer()
		log.Println("✅ Application server started!")
	}()

	// Start the UI
	log.Println("🛠️ Starting UI...")
	ui.RunUI()
	log.Println("✅ UI started!")
}
