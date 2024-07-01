package msg

type MsgType int
type Response int

const (
	MsgErr MsgType = iota
	Event
	Echo
	Exit

	ResErr Response = iota
	Ok
)
