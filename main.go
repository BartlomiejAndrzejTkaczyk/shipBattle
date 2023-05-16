package main

import (
	"log"
	"shipBattle/api"
	"shipBattle/aplication"
	"shipBattle/screen"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := api.CreateClient()
	if err != nil {
		log.Fatal(err)
	}

	s := screen.DefaultScreen()
	app := aplication.App{
		Client: client,
		Screen: s,
	}

	app.Run()
}
