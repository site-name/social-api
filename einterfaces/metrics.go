package einterfaces

import (
	"github.com/mattermost/logr/v2"
)

type MetricsInterface interface {
	Register()

	IncrementPostCreate()
	IncrementWebhookPost()
	IncrementPostSentEmail()
	IncrementPostSentPush()
	IncrementPostBroadcast()
	IncrementPostFileAttachment(count int)

	IncrementHttpRequest()
	IncrementHttpError()

	IncrementClusterRequest()
	ObserveClusterRequestDuration(elapsed float64)
	IncrementClusterEventType(eventType string)

	IncrementLogin()
	IncrementLoginFail()

	IncrementEtagHitCounter(route string)
	IncrementEtagMissCounter(route string)

	IncrementMemCacheHitCounter(cacheName string)
	IncrementMemCacheMissCounter(cacheName string)
	IncrementMemCacheInvalidationCounter(cacheName string)
	IncrementMemCacheMissCounterSession()
	IncrementMemCacheHitCounterSession()
	IncrementMemCacheInvalidationCounterSession()

	IncrementWebsocketEvent(eventType string)
	IncrementWebSocketBroadcast(eventType string)
	IncrementWebSocketBroadcastBufferSize(hub string, amount float64)
	DecrementWebSocketBroadcastBufferSize(hub string, amount float64)
	IncrementWebSocketBroadcastUsersRegistered(hub string, amount float64)
	DecrementWebSocketBroadcastUsersRegistered(hub string, amount float64)

	AddMemCacheHitCounter(cacheName string, amount float64)
	AddMemCacheMissCounter(cacheName string, amount float64)

	IncrementPostsSearchCounter()
	ObservePostsSearchDuration(elapsed float64)
	IncrementFilesSearchCounter()
	ObserveFilesSearchDuration(elapsed float64)
	ObserveStoreMethodDuration(method, success string, elapsed float64)
	ObserveApiEndpointDuration(endpoint, method, statusCode string, elapsed float64)
	IncrementPostIndexCounter()
	IncrementFileIndexCounter()
	IncrementUserIndexCounter()
	IncrementChannelIndexCounter()

	ObservePluginHookDuration(pluginID, hookName string, success bool, elapsed float64)
	ObservePluginMultiHookIterationDuration(pluginID string, elapsed float64)
	ObservePluginMultiHookDuration(elapsed float64)
	ObservePluginApiDuration(pluginID, apiName string, success bool, elapsed float64)

	ObserveEnabledUsers(users int64)
	GetLoggerMetricsCollector() logr.MetricsCollector

	IncrementRemoteClusterMsgSentCounter(remoteID string)
	IncrementRemoteClusterMsgReceivedCounter(remoteID string)
	IncrementRemoteClusterMsgErrorsCounter(remoteID string, timeout bool)
	ObserveRemoteClusterPingDuration(remoteID string, elapsed float64)
	ObserveRemoteClusterClockSkew(remoteID string, skew float64)
	IncrementRemoteClusterConnStateChangeCounter(remoteID string, online bool)

	IncrementJobActive(jobType string)
	DecrementJobActive(jobType string)

	SetReplicaLagAbsolute(node string, value float64)
	SetReplicaLagTime(node string, value float64)
}
