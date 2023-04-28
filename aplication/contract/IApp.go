package contract

import (
	"shipBattle/api"
	"shipBattle/screen"
)

//// todo temp, remove when real response will be ready
//type StatusResponse struct {
//}

type IApp interface {
	Run()
	initGame(client api.Client, screen screen.Screen) error
	Board() ([]string, error)
	Status() (*api.StatusResponse, error)
	Fire(coord string) (string, error)
}
