package eventsource

// import (
// 	"context"
// 	"time"

// 	"github.com/sitename/sitename/modules/log"
// 	"github.com/sitename/sitename/modules/timeutil"
// 	"github.com/sitename/sitename/models"
// 	"github.com/sitename/sitename/modules/graceful"
// 	"github.com/sitename/sitename/modules/setting"
// )

// // Init starts this eventsource
// func (m *Manager) Init() {
// 	if setting.UI.Notification.EventSourceUpdateTime <= 0 {
// 		return
// 	}
// 	go graceful.GetManager().RunWithShutdownContext(m.Run)
// }

// // Run runs the manager within a provided context
// func (m *Manager) Run(ctx context.Context) {
// 	then := timeutil.TimeStampNow().Add(-2)
// 	timer := time.NewTicker(setting.UI.Notification.EventSourceUpdateTime)

// loop:
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			timer.Stop()
// 			break loop
// 		case <-timer.C:
// 			now := timeutil.TimeStampNow().Add(-2)

// 			uidCounts, err := models.GetUIDsAndNotificationCounts(then, now)
// 			if err != nil {
// 				log.Error("Unable to get UIDcounts: %v", err)
// 			}
// 			for _, uidCount := range uidCounts {
// 				m.SendMessage(uidCount.UserID, &Event{
// 					Name: "notification-count",
// 					Data: uidCount,
// 				})
// 			}
// 			then = now
// 		}
// 	}
// 	m.UnregisterAll()
// }
