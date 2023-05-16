package screen

import (
	"bufio"
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"os"
	modelsApi "shipBattle/api/models"
	"shipBattle/screen/models"
	"strconv"
	"strings"
	"time"
)

type Screen struct {
	ui              *gui.GUI
	myBoard         *gui.Board
	opponenBoard    *gui.Board
	myNick          *gui.Text
	opponentNick    *gui.Text
	coordTxt        *gui.Text
	fireInfo        *gui.Text
	myStatus        [10][10]gui.State
	oppoStatus      [10][10]gui.State
	exitKayTxt      *gui.Text
	allDrawElements []gui.Drawable
}

func DefaultScreen() *Screen {
	return &Screen{
		ui:           gui.NewGUI(true),
		myStatus:     createEmptyStats(),
		oppoStatus:   createEmptyStats(),
		myBoard:      gui.NewBoard(3, 6, nil),
		opponenBoard: gui.NewBoard(50, 6, nil),
		fireInfo:     gui.NewText(3, 2, "", nil),
		exitKayTxt:   gui.NewText(3, 4, "Press CTRL+C to exit", nil),
	}
}

func createEmptyStats() [10][10]gui.State {
	state := [10][10]gui.State{}

	for i := 0; i < 10; i++ {
		state[i] = [10]gui.State{}
		for j := 0; j < 10; j++ {
			state[i][j] = gui.Empty
		}
	}

	return state
}

func (s *Screen) Start() {
	go s.ui.Start(nil)
}

func (s *Screen) ShowBattleFiled() {

	s.ui.Draw(s.myBoard)
	s.ui.Draw(s.opponenBoard)
	s.ui.Draw(s.exitKayTxt)
	s.ui.Draw(s.fireInfo)
}

func (s *Screen) SetMyShips(Positions []string) {
	for _, pos := range Positions {
		col, row := getPosFromCoord(pos)
		s.myStatus[col][row] = gui.Ship
	}

	s.myBoard.SetStates(s.myStatus)
}

func (s *Screen) SetNicknames(myNick string, opponentNick string) {
	if opponentNick == "" {
		opponentNick = "none opponent nickname"
	}

	s.myNick = gui.NewText(3, 28, myNick, nil)
	s.opponentNick = gui.NewText(50, 28, opponentNick, nil)

	s.ui.Draw(s.myNick)
	s.ui.Draw(s.opponentNick)
}

func (s *Screen) DrawShouldFire(ctx context.Context) {
	timeout, _ := ctx.Deadline()
	min := timeout.Minute() - time.Now().Minute()
	sec := timeout.Second() - time.Now().Second() + min*60

	for ; sec > 0; sec-- {
		s.fireInfo.SetBgColor(gui.Green)
		s.fireInfo.SetText("Choose your target (" + strconv.Itoa(sec) + " sec)")

		select {
		case <-ctx.Done():
			s.fireInfo.SetBgColor(gui.Red)
			s.fireInfo.SetText("Opponent's turn")
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (s *Screen) GetFirePos() *string {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	go s.DrawShouldFire(ctxTimeout)

	for true {
		pos := s.opponenBoard.Listen(ctxTimeout)
		if pos == "" {
			return nil
		}

		col, row := getPosFromCoord(pos)
		if s.oppoStatus[col][row] != gui.Empty {
			continue
		}
		return &pos
	}
	return nil
}

func (s *Screen) SetOppoShots(shots []string) {
	for _, shot := range shots {
		col, row := getPosFromCoord(shot)
		status := s.myStatus[col][row]

		if status == gui.Ship {
			s.myStatus[col][row] = gui.Hit
		} else if status == gui.Empty {
			s.myStatus[col][row] = gui.Miss
		}
	}

	s.myBoard.SetStates(s.myStatus)
}

func (s *Screen) SetMyShot(shotPos *string, isHit bool) {
	col, row := getPosFromCoord(*shotPos)
	newStatus := gui.Miss

	if isHit {
		newStatus = gui.Hit
	}

	s.oppoStatus[col][row] = newStatus
	s.opponenBoard.SetStates(s.oppoStatus)
}

func getPosFromCoord(coord string) (int, int) {
	col := coord[0] - 'A'
	row, _ := strconv.Atoi(coord[1:])
	return int(col), row - 1
}

func (s *Screen) Log(format string, args ...interface{}) {
	s.ui.Log(format, args...)
}

func (s *Screen) UserMenu(opponents []modelsApi.OpponentInfo) *models.GameOptions {
	gameOptions := models.GameOptions{}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Do you want to play with a bot? (y/n): ")
	isBot, err := reader.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}
	isBot = strings.TrimSpace(isBot)
	gameOptions.IsBot = isBot == "y"

	if !gameOptions.IsBot {
		fmt.Println("Select an opponent:")
		for i, oppo := range opponents {
			fmt.Printf("%d. %s (%s)\n", i+1, oppo.Nick, oppo.GameStatus)
		}
		fmt.Print("Enter the number of your opponent (0 for none): ")
		oppoNumStr, err := reader.ReadString('\n')
		if err != nil {
			log.Panic(err)
		}
		oppoNumStr = strings.TrimSpace(oppoNumStr)
		oppoNum, err := strconv.Atoi(oppoNumStr)
		if err != nil || oppoNum < 0 || oppoNum > len(opponents) {
			fmt.Println("Invalid input. Please enter a valid opponent number.")
			return s.UserMenu(opponents)
		}
		if oppoNum == 0 {
			gameOptions.OpponentNick = ""
		} else {
			gameOptions.OpponentNick = opponents[oppoNum-1].Nick
		}
	}
	fmt.Print("Enter your nickname: ")
	playerNick, err := reader.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}
	playerNick = strings.TrimSpace(playerNick)
	gameOptions.PlayerNick = playerNick
	return &gameOptions
}

func (s *Screen) EndScreen(status *modelsApi.StatusResponse) {
	s.ui.Remove(s.myBoard)
	s.ui.Remove(s.opponenBoard)
	s.ui.Remove(s.exitKayTxt)
	s.ui.Remove(s.fireInfo)

	s.ui.Draw(gui.NewText(0, 0, status.Desc, nil))
}
