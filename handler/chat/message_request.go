package chat

type MessageRequest struct {
	ChatType  byte
	ChatIndex byte
	Receiver  string
	Message   string
}
