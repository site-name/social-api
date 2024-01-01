package testlib

import (
	"sync"

	"github.com/sitename/sitename/einterfaces"
	"github.com/sitename/sitename/model_helper"
)

type FakeClusterInterface struct {
	clusterMessageHandler einterfaces.ClusterMessageHandler
	mut                   sync.RWMutex
	messages              []*model_helper.ClusterMessage
}

func (c *FakeClusterInterface) StartInterNodeCommunication() {}

func (c *FakeClusterInterface) StopInterNodeCommunication() {}

func (c *FakeClusterInterface) RegisterClusterMessageHandler(event model_helper.ClusterEvent, crm einterfaces.ClusterMessageHandler) {
	c.clusterMessageHandler = crm
}

func (c *FakeClusterInterface) HealthScore() int {
	return 0
}

func (c *FakeClusterInterface) GetClusterId() string { return "" }

func (c *FakeClusterInterface) IsLeader() bool { return false }

func (c *FakeClusterInterface) GetMyClusterInfo() *model_helper.ClusterInfo { return nil }

func (c *FakeClusterInterface) GetClusterInfos() []*model_helper.ClusterInfo { return nil }

func (c *FakeClusterInterface) SendClusterMessage(message *model_helper.ClusterMessage) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.messages = append(c.messages, message)
}

func (c *FakeClusterInterface) SendClusterMessageToNode(nodeID string, message *model_helper.ClusterMessage) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.messages = append(c.messages, message)
	return nil
}

func (c *FakeClusterInterface) NotifyMsg(buf []byte) {}

func (c *FakeClusterInterface) GetClusterStats() ([]*model_helper.ClusterStats, *model_helper.AppError) {
	return nil, nil
}

func (c *FakeClusterInterface) GetLogs(page, perPage int) ([]string, *model_helper.AppError) {
	return []string{}, nil
}

func (c *FakeClusterInterface) QueryLogs(page, perPage int) (map[string][]string, *model_helper.AppError) {
	return make(map[string][]string), nil
}

func (c *FakeClusterInterface) ConfigChanged(previousConfig *model_helper.Config, newConfig *model_helper.Config, sendToOtherServer bool) *model_helper.AppError {
	return nil
}

func (c *FakeClusterInterface) SendClearRoleCacheMessage() {
	if c.clusterMessageHandler != nil {
		c.clusterMessageHandler(&model_helper.ClusterMessage{
			Event: model_helper.ClusterEventInvalidateCacheForRoles,
		})
	}
}

func (c *FakeClusterInterface) GetPluginStatuses() (model_helper.PluginStatuses, *model_helper.AppError) {
	return nil, nil
}

func (c *FakeClusterInterface) GetMessages() []*model_helper.ClusterMessage {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.messages
}

func (c *FakeClusterInterface) SelectMessages(filterCond func(message *model_helper.ClusterMessage) bool) []*model_helper.ClusterMessage {
	c.mut.RLock()
	defer c.mut.RUnlock()

	filteredMessages := []*model_helper.ClusterMessage{}
	for _, msg := range c.messages {
		if filterCond(msg) {
			filteredMessages = append(filteredMessages, msg)
		}
	}
	return filteredMessages
}

func (c *FakeClusterInterface) ClearMessages() {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.messages = nil
}
