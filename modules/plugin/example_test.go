package plugin_test

import (
	"fmt"
	"net/http"

	"github.com/sitename/sitename/modules/plugin"
)

type HelloWorldPlugin struct {
	plugin.SitenamePlugin
}

func (p *HelloWorldPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func Example_helloWorld() {
	plugin.ClientMain(&HelloWorldPlugin{})
}
