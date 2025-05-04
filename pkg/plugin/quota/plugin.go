package quota

import (
	log "log/slog"

	"github.com/policyd/pkg/types"
	"gopkg.in/yaml.v3"
)

// config configuration to Quota plugin
type config struct {
	Sql string `yaml:"sql"`
}

func (c *config) Load(n yaml.Node) error {
	return n.Decode(c)
}

type Plugin struct {
}

func (q *Plugin) Config(m map[string]yaml.Node) error {
	c := config{}
	if err := c.Load(m[q.Name()]); err != nil {
		return err
	}
	return nil
}

func (p *Plugin) Name() string {
	return "quota"
}

func (p *Plugin) PostMessageReceived(types.RequestID, types.Message) {
	log.Info("quota", "action", "PostMessage")

}

func (q *Plugin) PreResponse(types.RequestID, *types.Response) *types.Response {
	log.Info("quota", "action", "PreMessage")
	return types.OkResponse
}
