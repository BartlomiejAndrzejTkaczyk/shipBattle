package main

import (
	"log"
	"shipBattle/api"
	"shipBattle/aplication"
	"shipBattle/screen"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := aplication.App{}
	client, _ := api.CreateClient()
	app.Run(client, screen.DefaultScreen())
}
