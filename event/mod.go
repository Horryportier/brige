package event

type EventRegistry map[string][]string

type EventDescription struct {
	Name string
	Args map[string]any
}


