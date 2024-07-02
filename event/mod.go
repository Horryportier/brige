package event

type Event struct {
	Name string `json:"name"`
	Data string `json:"cmd"`
}

type EventTriger struct {
	Name string
	Cmd  string
}
