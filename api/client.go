package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	modelsApi "shipBattle/api/models"
	"shipBattle/configuration"
	"shipBattle/screen/models"
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

func (c *Client) StartGameConnect(gameOptions *models.GameOptions) {
	log.Print("Connecting ...")
	postBody, _ := json.Marshal(struct {
		Nick       string `json:"nick"`
		TargetNick string `json:"target_nick"`
		Wpbot      bool   `json:"wpbot"`
	}{
		Nick:       *gameOptions.PlayerNick,
		TargetNick: *gameOptions.OpponentNick,
		Wpbot:      *gameOptions.IsBot,
	})

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(c.apiCfg.InitUrl, "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	c.token = resp.Header.Get(c.apiCfg.AuthTokenName)
	log.Print("Connected")
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

func (c *Client) Status(ctx context.Context) (*modelsApi.StatusResponse, error) {
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
	var body modelsApi.StatusResponse
	err = json.Unmarshal(bodyByte, &body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	log.Printf("%+v\n", body)
	return &body, nil
}

func (c *Client) Fire(shot *string) bool {
	bodyReq, _ := json.Marshal(modelsApi.FireReq{
		Coord: *shot,
	})

	bodyResp, err := c.request(context.Background(), http.MethodPost, c.apiCfg.FireUrl, bytes.NewReader(bodyReq))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	var fireResp modelsApi.FireResp
	err = json.Unmarshal(bodyResp, &fireResp)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	return fireResp.Result == "hit"
}

func (c *Client) OppoList() []modelsApi.OpponentInfo {
	bodyResp, err := c.request(context.Background(), http.MethodGet, c.apiCfg.OppoListUrl, nil)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	if len(bodyResp) == 4 {
		return []modelsApi.OpponentInfo{}
	}

	var oppoList []modelsApi.OpponentInfo
	err = json.Unmarshal(bodyResp, &oppoList)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	return oppoList
}

func (c *Client) request(ctx context.Context, method string, url string, bodyReq *bytes.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
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
	return bodyByte, err
}
