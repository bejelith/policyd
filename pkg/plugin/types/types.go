package types

import (
	"github.com/policyd/pkg/types"
	"gopkg.in/yaml.v3"
)

type Plugin interface {
	Name() string
	Config(map[string]yaml.Node) error
	PostMessageReceived(types.RequestID, types.Message)
	PreResponse(types.RequestID, *types.Response) *types.Response
}
