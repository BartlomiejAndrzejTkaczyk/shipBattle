package aplication

import (
	"context"
	"log"
	"shipBattle/api"
	"shipBattle/screen"
	"time"
)

type App struct {
	client *api.Client
	screen *screen.Screen
}

func (app *App) Run(client *api.Client, screen *screen.Screen) {
	app.client = client
	app.screen = screen
	ticker := time.NewTimer(1 * time.Second)

	app.client.Connect()
	for _ = range ticker.C {
		status, _ := app.Status()
		if status.GameStatus != "waiting" {
			app.screen.SetNickname(status.Nick)
			break
		}
		ticker.Reset(1 * time.Second)
	}
	go app.screen.Start()
	app.Board()
	time.Sleep(5 * time.Second)
}

func (app *App) Board() {
	coords, err := app.client.GetBoard(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	app.screen.SetMyShips(coords)
}

func (app *App) Status() (*api.StatusResponse, error) {
	status, _ := app.client.Status(context.Background())
	return status, nil
}
