package main

import(
	"fmt"
	"github.com/ArteShow/Calculator/application"
	"github.com/ArteShow/Calculator/pkg/setup"

)

func main(){
	fmt.Println("Starting application...")
	setup.Setup()
	application.StartApplicationServer()
}
