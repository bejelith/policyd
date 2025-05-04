package handler

import (
	log "log/slog"

	"github.com/policyd/pkg/plugin"

	"github.com/policyd/pkg/types"
)

var defaultMessageHandler = handleMessage

func DefaultMessageHander() MessageHandler {
	return defaultMessageHandler
}

type MessageHandler func(types.RequestID, types.Message) *types.Response

func handleMessage(id types.RequestID, message types.Message) *types.Response {
	log.Debug("received", "message", message)
	for _, p := range plugin.Enumerate() {
		p.PostMessageReceived(id, message)
	}
	response := types.OkResponse
	for _, p := range plugin.Enumerate() {
		response = p.PreResponse(id, response)
		if response.Resp != types.OkResponse.Resp {
			return response
		}
	}
	return response
}
