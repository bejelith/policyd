package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type quota struct {
	Sql string `yaml:"sql"`
}

var validYaml = `---
host: localhost
port: 12345
plugins:
  quota:
    sql: localhost
  http:
    sql: sql
`
var invalidYaml = `---
host: localhost
port: asd
plugins:
  quota:
    sql: localhost
  http:
    sql: sql
`

var noPluginYaml = `---
host: localhost
port: 12345
plugins:
  
`

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		yml       string
		hasPlugin bool
		err       error
	}{{
		"ValidYaml",
		validYaml,
		true,
		nil,
	}, {
		"InvalidYaml",
		invalidYaml,
		true,
		&yaml.TypeError{},
	}, {
		"NoPlugins",
		noPluginYaml,
		false,
		nil,
	}}
	for _, test := range tests {
		t.Run(test.name, func(tb *testing.T) {
			f, err := os.CreateTemp("", "*")
			tb.Cleanup(func() { os.Remove(f.Name()) })
			if err != nil {
				t.Fatal(err)
			}
			if _, err := f.WriteString(test.yml); err != nil {
				t.Fatal(err)
			}
			c, e := ParseFile(f.Name())
			assert.IsType(t, test.err, e, e)
			if e != nil {
				return
			}
			assert.Equal(t, "localhost", c.Host)
			assert.Equal(t, 12345, c.Port)

			plugins := map[string]yaml.Node{}
			assert.Nil(t, c.Plugins.Decode(plugins))
			if test.hasPlugin {
				quotaPlugin := &quota{}
				r := plugins["quota"]
				assert.Nil(t, r.Decode(quotaPlugin))
				assert.Equal(t, "localhost", quotaPlugin.Sql)
			}
		})
	}
}

func TestReadFromFileError(t *testing.T) {
	_, e := ParseFile("nofile")
	assert.IsType(t, &os.PathError{}, e)
}
