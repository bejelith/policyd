package quota_plugin

import "githun.com/policyd/pkg/plugin"

type quota_plugin struct {
}

func (q *quota_plugin) PostMessageReceived(plugin.RequestID, []string) {

}

func (q *quota_plugin) PreResponse(plugin.RequestID, *plugin.Response) *plugin.Response {
	return plugin.OkResponse
}

func init() {
	plugin.Register(&quota_plugin{})
}
