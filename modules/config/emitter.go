package config

import (
	"sync"

	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
)

// Listener is a callback function invoked when the configuration changes.
type Listener func(oldCfg, newCfg *model_helper.Config)

// emitter enables threadsafe registration and broadcasting to configuration listeners
type emitter struct {
	listeners sync.Map
}

// AddListener adds a callback function to invoke when the configuration is modified.
func (e *emitter) AddListener(listener Listener) string {
	id := model_helper.NewId()
	e.listeners.Store(id, listener)
	return id
}

// RemoveListener removes a callback function using an id returned from AddListener.
func (e *emitter) RemoveListener(id string) {
	e.listeners.Delete(id)
}

// invokeConfigListeners synchronously notifies all listeners about the configuration change.
func (e *emitter) invokeConfigListeners(oldCfg, newCfg *model_helper.Config) {
	e.listeners.Range(func(key, value any) bool {
		listener := value.(Listener)
		listener(oldCfg, newCfg)
		return true
	})
}

// srcEmitter enables threadsafe registration and broadcasting to configuration listeners
type logSrcEmitter struct {
	listeners sync.Map
}

// AddListener adds a callback function to invoke when the configuration is modified.
func (e *logSrcEmitter) AddListener(listener LogSrcListener) string {
	id := model_helper.NewId()
	e.listeners.Store(id, listener)
	return id
}

// RemoveListener removes a callback function using an id returned from AddListener.
func (e *logSrcEmitter) RemoveListener(id string) {
	e.listeners.Delete(id)
}

// invokeConfigListeners synchronously notifies all listeners about the configuration change.
func (e *logSrcEmitter) invokeConfigListeners(oldCfg, newCfg slog.LoggerConfiguration) {
	e.listeners.Range(func(key, value any) bool {
		listener := value.(LogSrcListener)
		listener(oldCfg, newCfg)
		return true
	})
}
