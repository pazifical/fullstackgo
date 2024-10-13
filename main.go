package main

import (
	"html/template"
	"log"
	"os"

	"github.com/pazifical/fullstackgo/internal"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Please provide either 'dev' oder 'build' mode command line argument")
	}
	mode := os.Args[1]

	if mode == "build" {
		log.Println("Initializing Build Mode")
		err := startBuildMode()
		if err != nil {
			log.Fatal(err)
		}
	} else if mode == "dev" {
		log.Println("Initializing Developer Mode")
		err := startDevMode()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func startDevMode() error {
	developerMode := internal.NewDeveloperMode()
	return developerMode.Start()
}

func startBuildMode() error {
	t, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return err
	}

	generator := internal.NewGenerator(t)
	return generator.Run()
}
