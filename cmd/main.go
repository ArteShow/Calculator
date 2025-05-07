package main

import(
	"fmt"
	"github.com/ArteShow/Calculator/application"
	"github.com/ArteShow/Calculator/pkg/setup"
	"github.com/ArteShow/Calculator/internal"


)

func main(){
	fmt.Println("Starting application...")
	setup.Setup()
	go func(){
		internal.StartTPCListener()
	}()
	application.StartApplicationServer()
}
