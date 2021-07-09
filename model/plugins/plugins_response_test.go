package plugins

import (
	"strings"
	"testing"

	"github.com/sitename/sitename/model"
	"github.com/stretchr/testify/assert"
)

func TestPluginsResponseJson(t *testing.T) {
	manifest := &model.Manifest{
		Id: "theid",
		Server: &model.ManifestServer{
			Executable: "theexecutable",
		},
		Webapp: &model.ManifestWebapp{
			BundlePath: "thebundlepath",
		},
	}

	response := &PluginsResponse{
		Active:   []*PluginInfo{{Manifest: *manifest}},
		Inactive: []*PluginInfo{},
	}

	json := response.ToJson()
	newResponse := PluginsResponseFromJson(strings.NewReader(json))
	assert.Equal(t, newResponse, response)
	assert.Equal(t, newResponse.ToJson(), json)
	assert.Equal(t, PluginsResponseFromJson(strings.NewReader("junk")), (*PluginsResponse)(nil))
}
