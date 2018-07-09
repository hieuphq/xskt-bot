package bot

// Message .
type Message interface{}

// Bot .
type Bot interface {
	ReceiveMessage() (Message, error)
	SendMessage(message Message) error
}
