package screen

import (
	gui "github.com/grupawp/warships-gui/v2"
	"strconv"
)

type Screen struct {
	ui           *gui.GUI
	myBoard      *gui.Board
	opponenBoard *gui.Board
	myNick       *gui.Text
	opponentNick *gui.Text
	coordTxt     *gui.Text
	//todo exitKayTxt *gui.Text
}

func DefaultScreen() *Screen {
	return &Screen{
		ui: gui.NewGUI(true),
	}
}

func CreateEmptyStats() [10][10]gui.State {
	state := [10][10]gui.State{}
	for i := 0; i < 10; i++ {
		state[i] = [10]gui.State{}
	}
	return state
}

func (s *Screen) Init() {
	s.myBoard = gui.NewBoard(3, 6, nil)
	s.opponenBoard = gui.NewBoard(50, 6, nil)
	s.opponentNick = gui.NewText(3, 19, "mynick", nil)
	s.opponentNick = gui.NewText(50, 19, "opponentnick", nil)

	s.ui.Draw(gui.NewText(3, 2, "[coord]", nil))
	s.ui.Draw(gui.NewText(3, 4, "Press CTRL+C to exit", nil))
	s.ui.Draw(s.myNick)
	s.ui.Draw(s.opponentNick)
	s.ui.Draw(s.myBoard)
	s.ui.Draw(s.opponenBoard)
}

func (s *Screen) SetMyShips(coords []string) {
	myState := CreateEmptyStats()
	for _, coord := range coords {
		col := coord[0] - 'A'
		row, _ := strconv.Atoi(coord[1:])
		myState[col][row-1] = gui.Ship
	}
	s.myBoard.SetStates(myState)
	s.ui.Draw(s.myBoard)
}

func (s *Screen) SetNickname(nickname string) {
	s.myNick = gui.NewText(3, 28, nickname, nil)
	s.ui.Draw(s.myNick)
}

func (s *Screen) Log(msg any) {
	s.ui.Log("msg", msg)
}

func (s *Screen) Start() {
	s.Init()
	s.ui.Start(nil)
}
