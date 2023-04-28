package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"shipBattle/configuration"
)

type Client struct {
	token  string
	apiCfg *configuration.ApiCfg
}

func CreateClient() (*Client, error) {
	apiCfg, err := configuration.GetApiConfig()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c := &Client{
		apiCfg: apiCfg,
	}

	return c, nil
}

type connectReq struct {
	Coords     []string `json:"coords"`
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick"`
	Wpbot      bool     `json:"wpbot"`
}

func (c *Client) Connect() {
	log.Print("Connecting ...")
	postBody, _ := json.Marshal(connectReq{
		Wpbot: true,
	})

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(c.apiCfg.InitUrl, "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	log.Print("Connected")
	c.token = resp.Header.Get(c.apiCfg.AuthTokenName)
}

func (c *Client) GetBoard(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiCfg.BoardUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(c.apiCfg.AuthTokenName, c.token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyByte, err := io.ReadAll(resp.Body)
	var body map[string][]string
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	board := body["board"]
	return board, nil
}

type StatusResponse struct {
	Desc           string   `json:"desc"`
	GameStatus     string   `json:"game_status"`
	LastGameStatus string   `json:"last_game_status"`
	Nick           string   `json:"nick"`
	OppDesc        string   `json:"opp_desc"`
	OppShots       []string `json:"opp_shots"`
	Opponent       string   `json:"opponent"`
	ShouldFire     bool     `json:"should_fire"`
	Timer          int      `json:"timer"`
}

func (c *Client) Status(ctx context.Context) (*StatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://go-pjatk-server.fly.dev/api/game", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(c.apiCfg.AuthTokenName, c.token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyByte, err := io.ReadAll(resp.Body)
	var body StatusResponse
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	return &body, nil
}
