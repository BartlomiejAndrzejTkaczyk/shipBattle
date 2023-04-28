package contract

type IScreen interface {
	Init()
	SetNickname(nickname string)
	Log(msg any)
}
