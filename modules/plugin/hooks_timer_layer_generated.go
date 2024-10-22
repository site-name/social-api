// Code generated by "make pluginapi"
// DO NOT EDIT

package plugin

import (
	"io"
	"net/http"
	timePkg "time"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
)

type hooksTimerLayer struct {
	pluginID  string
	hooksImpl Hooks
	metrics   einterfaces.MetricsInterface
}

func (hooks *hooksTimerLayer) recordTime(startTime timePkg.Time, name string, success bool) {
	if hooks.metrics != nil {
		elapsedTime := float64(timePkg.Since(startTime)) / float64(timePkg.Second)
		hooks.metrics.ObservePluginHookDuration(hooks.pluginID, name, success, elapsedTime)
	}
}

func (hooks *hooksTimerLayer) OnActivate() error {
	startTime := timePkg.Now()
	_returnsA := hooks.hooksImpl.OnActivate()
	hooks.recordTime(startTime, "OnActivate", _returnsA == nil)
	return _returnsA
}

func (hooks *hooksTimerLayer) Implemented() ([]string, error) {
	startTime := timePkg.Now()
	_returnsA, _returnsB := hooks.hooksImpl.Implemented()
	hooks.recordTime(startTime, "Implemented", _returnsB == nil)
	return _returnsA, _returnsB
}

func (hooks *hooksTimerLayer) OnDeactivate() error {
	startTime := timePkg.Now()
	_returnsA := hooks.hooksImpl.OnDeactivate()
	hooks.recordTime(startTime, "OnDeactivate", _returnsA == nil)
	return _returnsA
}

func (hooks *hooksTimerLayer) OnConfigurationChange() error {
	startTime := timePkg.Now()
	_returnsA := hooks.hooksImpl.OnConfigurationChange()
	hooks.recordTime(startTime, "OnConfigurationChange", _returnsA == nil)
	return _returnsA
}

func (hooks *hooksTimerLayer) ServeHTTP(c *Context, w http.ResponseWriter, r *http.Request) {
	startTime := timePkg.Now()
	hooks.hooksImpl.ServeHTTP(c, w, r)
	hooks.recordTime(startTime, "ServeHTTP", true)
}

func (hooks *hooksTimerLayer) UserHasBeenCreated(c *Context, user *model.User) {
	startTime := timePkg.Now()
	hooks.hooksImpl.UserHasBeenCreated(c, user)
	hooks.recordTime(startTime, "UserHasBeenCreated", true)
}

func (hooks *hooksTimerLayer) UserWillLogIn(c *Context, user *model.User) string {
	startTime := timePkg.Now()
	_returnsA := hooks.hooksImpl.UserWillLogIn(c, user)
	hooks.recordTime(startTime, "UserWillLogIn", true)
	return _returnsA
}

func (hooks *hooksTimerLayer) UserHasLoggedIn(c *Context, user *model.User) {
	startTime := timePkg.Now()
	hooks.hooksImpl.UserHasLoggedIn(c, user)
	hooks.recordTime(startTime, "UserHasLoggedIn", true)
}

func (hooks *hooksTimerLayer) FileWillBeUploaded(c *Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {
	startTime := timePkg.Now()
	_returnsA, _returnsB := hooks.hooksImpl.FileWillBeUploaded(c, info, file, output)
	hooks.recordTime(startTime, "FileWillBeUploaded", true)
	return _returnsA, _returnsB
}

func (hooks *hooksTimerLayer) OnPluginClusterEvent(c *Context, ev model_helper.PluginClusterEvent) {
	startTime := timePkg.Now()
	hooks.hooksImpl.OnPluginClusterEvent(c, ev)
	hooks.recordTime(startTime, "OnPluginClusterEvent", true)
}
