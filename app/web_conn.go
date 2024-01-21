package app

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/i18n"
)

const (
	sendQueueSize          = 256
	sendSlowWarn           = (sendQueueSize * 50) / 100
	sendFullWarn           = (sendQueueSize * 95) / 100
	writeWaitTime          = 30 * time.Second
	pongWaitTime           = 100 * time.Second
	pingInterval           = (pongWaitTime * 6) / 10
	authCheckInterval      = 5 * time.Second
	webConnMemberCacheTime = 1000 * 60 * 30 // 30 minutes
	deadQueueSize          = 128            // Approximated from /proc/sys/net/core/wmem_default / 2048 (avg msg size)
)

const (
	reconnectFound    = "success"
	reconnectNotFound = "failure"
	reconnectLossless = "lossless"
)

type WebConnConfig struct {
	WebSocket    *websocket.Conn
	Session      model.Session
	TFunc        i18n.TranslateFunc
	Locale       string
	ConnectionID string
	Active       bool

	// unexported
	sequence         int
	activeQueue      chan model_helper.WebSocketMessage
	deadQueue        []*model_helper.WebSocketEvent
	deadQueuePointer int
}

// WebConn represents a single websocket connection to a user.
// It contains all the necessary state to manage sending/receiving data to/from
// a websocket.
type WebConn struct {
	sessionExpiresAt int64
}
