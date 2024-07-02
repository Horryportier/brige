package msg

type MsgType int
type Response int

const (
	MsgErr MsgType = iota
	Inital
	Event
	Echo
	Exit

	ResErr Response = iota
	Ok
)
