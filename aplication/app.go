package aplication

import (
	"context"
	"log"
	"shipBattle/api"
	modelsApi "shipBattle/api/models"
	"shipBattle/screen"
	"shipBattle/screen/models"
	"time"
)

type App struct {
	Client      *api.Client
	Screen      *screen.Screen
	gameOptions *models.GameOptions
}

func (app *App) Run() {
	app.gameOptions = app.Screen.UserMenu(app.Client.OppoList())

	app.Screen.Start()
	app.beginBattle()
	app.mainBattleLoop()
	time.Sleep(4 * time.Second)
}

func (app *App) beginBattle() {
	app.Client.StartGameConnect(app.gameOptions)
	// Wait for game to start
	ticker := time.Tick(time.Second)
	for range ticker {
		status, _ := app.status()
		if status.GameStatus != "waiting" {

			app.Screen.SetNicknames(status.Nick, status.Opponent)
			app.board()
			return
		}
	}
}

func (app *App) mainBattleLoop() {
	app.Screen.ShowBattleFiled()

	//for {
	//	ctx, cancel := context.WithCancel(context.Background())
	//	statusCh := make(chan *modelsApi.StatusResponse)
	//	go app.waitToStatus(statusCh, cancel)
	//
	//	select {
	//	case <-ctx.Done():
	//		return
	//	case status := <-statusCh:
	//		switch status.GameStatus {
	//		case "game_in_progress":
	//			if status.ShouldFire {
	//				app.fire(status.OppShots)
	//			}
	//			cancel()
	//		case "waiting":
	//			time.Sleep(1 * time.Second)
	//			cancel()
	//		}
	//	}
	//}

	for {
		status, err := app.status()
		if err != nil {
			app.Screen.Log("Error: %s", err.Error())
		}
		switch status.GameStatus {
		case "game_in_progress":
			if status.ShouldFire {
				app.fire(status.OppShots)
			}
		case "ended":
			app.Screen.EndScreen(status)
		case "waiting":
		}
		time.Sleep(1 * time.Second)
	}

}

func (app *App) fire(oppoShots []string) {
	app.Screen.SetOppoShots(oppoShots)
	myShot := app.Screen.GetFirePos()
	if myShot != nil {
		isHit := app.Client.Fire(myShot)
		app.Screen.SetMyShot(myShot, isHit)
	}
}

func (app *App) board() {
	coords, err := app.Client.GetBoard(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	app.Screen.SetMyShips(coords)
}

func (app *App) status() (*modelsApi.StatusResponse, error) {
	status, err := app.Client.Status(context.Background())
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second)
	return status, err
}

func (app *App) waitToStatus(statusCh chan *modelsApi.StatusResponse, cancel context.CancelFunc) {
	status, err := app.status()

	if err != nil {
		app.Screen.Log("Error: ", err)
		cancel()
		return
	}

	statusCh <- status
}
