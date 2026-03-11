package entity

type Message struct {
	Name    string
	Headers []Header
	Content string
}

type MessageOutput struct {
	Name   string
	Status string
}

type Header struct {
	Key   string
	Value string
}
