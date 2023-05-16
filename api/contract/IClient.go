package contract

import (
	"context"
	"shipBattle/api"
)

type IClient interface {
	GetBoard(ctx context.Context) ([]string, error)
	StartGameConnect()
	Status(ctx context.Context) (*api.StatusResponse, error)
}
