package plugin

import (
	"sync"

	log "log/slog"

	"github.com/policyd/pkg/plugin/quota"
	"github.com/policyd/pkg/plugin/types"
	"gopkg.in/yaml.v3"
)

var plugins = []types.Plugin{
	&quota.Plugin{},
}
var lock = sync.RWMutex{}

func Configure(config yaml.Node) {
	pluginConfigs := map[string]yaml.Node{}

	if err := config.Decode(pluginConfigs); err != nil {
		log.Error("failed to parse plugin configuration", "error", err)
	}

	for _, plugin := range plugins {
		log.Info("registering plugin", "plugin", plugin.Name())
		lock.Lock()
		defer lock.Unlock()
		if err := plugin.Config(pluginConfigs); err != nil {
			log.Error("could not load plugin", "plugin", plugin.Name(), "error", err)
		}
	}
}

func Enumerate() []types.Plugin {
	lock.RLock()
	defer lock.RUnlock()
	n := make([]types.Plugin, len(plugins))
	copy(n, plugins)
	return n
}
