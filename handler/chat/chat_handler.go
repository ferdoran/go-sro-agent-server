package chat

import (
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/sirupsen/logrus"
)

const (
	All byte = iota + 1
	PM
	AllGM
	Party
	Guild
	Global
	Notice
	_
	Stall
	_
	Union
	_
	NPC
	_
	_
	Academy
)

type ChatHandler struct {
	channel chan server.PacketChannelData
}

func InitChatHandler() {
	queue := server.PacketManagerInstance.GetQueue(opcode.ChatRequest)
	handler := ChatHandler{channel: queue}
	go handler.Handle()
}

func (h *ChatHandler) Handle() {
	for {
		data := <-h.channel
		chatType, err := data.ReadByte()
		if err != nil {
			logrus.Panicln("Failed to read receiver")
		}
		chatIdx, err := data.ReadByte()
		if err != nil {
			logrus.Panicln("Failed to read receiver")
		}

		request := MessageRequest{
			ChatType:  chatType,
			ChatIndex: chatIdx,
			Receiver:  "",
			Message:   "",
		}

		if chatType == PM {
			receiver, err := data.ReadString()
			if err != nil {
				logrus.Panicln("Failed to read receiver")
			}
			logrus.Tracef("PM message receiver: %v\n", receiver)
			request.Receiver = receiver
		}

		msg, err := data.ReadString()
		if err != nil {
			logrus.Panicln("Failed to read message")
		}
		request.Message = msg

		switch request.ChatType {
		case All:
			handleAllMessage(request, data.Session)
		case Party:
			handlePartyMessage(request, data.Session)
		case Guild:

		case Global:
			handleGlobalMessage(request, data.Session)
		case Notice:

		case Stall:
			handleStallMessage(request, data.Session)
		case Union:

		case NPC:

		case Academy:

		case PM:
			handleWhisperMessage(request, data.Session)
		case AllGM:
			gmh := GetGmMessageHandlerInstance()
			gmh.HandleAdminMessage(request, data.Session)
		default:
			logrus.Debugf("unhandled chat message %v\n", request)
		}
	}
}
