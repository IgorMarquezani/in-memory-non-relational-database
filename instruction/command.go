package instruction

type Message struct {
	Err  error
	Data string
}

type Instruction struct {
	Command string
	Key     string 
	Data    any
	Type    string
	Channel chan Message
}
