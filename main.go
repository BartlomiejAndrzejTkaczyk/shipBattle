package main

import (
	"fmt"
	"log"
	"shipBattle/api"
	"shipBattle/aplication"
	"shipBattle/screen"
	"shipBattle/screen/models"
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

	var defaultMenu string
	fmt.Scan(&defaultMenu)

	if defaultMenu == "y" {
		isBot := false
		playerNick := "Player1"
		OpponentNick := "John_Doe_III"

		app.Run(
			models.GameOptions{IsBot: &isBot, PlayerNick: &playerNick, OpponentNick: &OpponentNick},
		)
	} else {
		app.Run(models.GameOptions{})
	}

}
