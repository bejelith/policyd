package handler

type MessageHandler interface {
	Handle([]string) string
}
