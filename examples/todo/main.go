package main

import (
	"log"

	"todo/app"

	grove "github.com/aryan-salemababdi/Grove/app"
)

func main() {
	a := grove.New()
	if err := a.RegisterModule(app.New()); err != nil {
		log.Fatal(err)
	}
	if err := a.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
