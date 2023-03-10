package plugin

import "sync"

var plugins = []Plugin{}
var lock = sync.RWMutex{}

type RequestID int64

type Plugin interface {
	PostMessageReceived(RequestID, []string)
	PreResponse(RequestID, *Response) *Response
}

type Response struct {
	Resp     string
	Failures []error
}

var OkResponse = &Response{Resp: "ok"}

func Register(plugin Plugin) {
	lock.Lock()
	defer lock.Unlock()
	plugins = append(plugins, plugin)
}

func Enumerate() []Plugin {
	lock.RLock()
	defer lock.RUnlock()
	n := make([]Plugin, len(plugins))
	copy(n, plugins)
	return n
}
