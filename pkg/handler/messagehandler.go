package handler

import "github.com/policyd/pkg/plugin"

type MessageHandler func(plugin.RequestID, []string) *plugin.Response

func handleMessage(id plugin.RequestID, message []string) *plugin.Response {
	for _, p := range plugin.Enumerate() {
		p.PostMessageReceived(id, message)
	}
	response := &plugin.Response{}
	for _, p := range plugin.Enumerate() {
		response = p.PreResponse(id, response)
	}
	return response
}
