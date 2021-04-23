package app

import (
	"sync"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/slog"
)

const (
	EmailBatchingTaskName = "Email Batching"
)

type batchedNotification struct {
	userID string
}

type EmailBatchingJob struct {
	server               *Server
	newNotifications     chan *batchedNotification
	pendingNotifications map[string][]*batchedNotification
	task                 *model.ScheduledTask
	taskMutex            sync.Mutex
}

func (es *EmailService) InitEmailBatching() {
	emailSetting := es.srv.Config().EmailSettings
	if *emailSetting.EnableEmailBatching {
		if es.EmailBatching == nil {
			es.EmailBatching = NewEmailBatchingJob(es, *emailSetting.EmailBatchingBufferSize)
		}

		// note that we don't support changing EmailBatchingBufferSize without restarting the server
		es.EmailBatching.Start()
	}
}

func NewEmailBatchingJob(es *EmailService, bufferSize int) *EmailBatchingJob {
	return &EmailBatchingJob{
		server:               es.srv,
		newNotifications:     make(chan *batchedNotification, bufferSize),
		pendingNotifications: make(map[string][]*batchedNotification),
	}
}

func (job *EmailBatchingJob) Start() {
	slog.Debug("Email batching job starting. Checking for pending emails periodically.", slog.Int("interval_in_seconds", *job.server.Config().EmailSettings.EmailBatchingInterval))
	newTask := model.CreateRecurringTask(EmailBatchingTaskName, job.CheckPendingEmails, time.Duration(*job.server.Config().EmailSettings.EmailBatchingInterval)*time.Second)

	job.taskMutex.Lock()
	oldTask := job.task
	job.task = newTask
	job.taskMutex.Unlock()

	if oldTask != nil {
		oldTask.Cancel()
	}
}

func (job *EmailBatchingJob) Add(user *model.User) bool {
	notification := &batchedNotification{userID: user.Id}
	select {
	case job.newNotifications <- notification:
		return true
	default:
		// return false if we couldn't queue the email notification so that we can send an immediate email
		return false
	}
}

func (job *EmailBatchingJob) CheckPendingEmails() {
	job.handleNewNotifications()

	// it's a bit weird to pass the send email function through here, but it makes it so that we can test
	// without actually sending emails
	job.checkPendingNotifications(time.Now(), job.server.EmailService.sendBatchedEmailNotification)
	slog.Debug("Email batching job ran. Some users still have notifications pending.", slog.Int("number_of_users", len(job.pendingNotifications)))
}

func (job *EmailBatchingJob) handleNewNotifications() {
	receiving := true

	// read in new notifications to send
	for receiving {
		select {
		case notification := <-job.newNotifications:
			userID := notification.userID

			if _, ok := job.pendingNotifications[userID]; !ok {
				job.pendingNotifications[userID] = []*batchedNotification{notification}
			} else {
				job.pendingNotifications[userID] = append(job.pendingNotifications[userID], notification)
			}
		default:
			receiving = false
		}
	}
}

func (job *EmailBatchingJob) checkPendingNotifications(now time.Time, handler func(string, []*batchedNotification)) {
	// for userID, notifications := range job.pendingNotifications {
	// 	// get how long we need to wait to send notifications to the user
	// 	preference, err := job.server.Store.Preference().Get(userID, model.PREFERENCE_CATEGORY_NOTIFICATIONS, model.PREFERENCE_NAME_EMAIL_INTERVAL)
	// 	if err != nil {
	// 		// use the default batching interval if an error ocurrs while fetching user preferences
	// 		interval, _ =
	// 	}
	// }
}

func (es *EmailService) sendBatchedEmailNotification(userID string, notifications []*batchedNotification) {
	// user, err := es.srv.Store.User().Get(context.Background(), userID)
	// if err != nil {
	// 	slog.Warn("Unable to find recipient for batched email notification")
	// 	return
	// }

	// translatedFunc := i18n.GetUserTranslations(user.Locale)
	// var content string
	// for _, notification := range notifications {
	// 	sender, err := es.srv.Store.User().Get(context.Background(), notification.userID)
	// 	if err != nil {
	// 		slog.Warn("Unable to find sender or post for batched email notification")
	// 		continue
	// 	}

	// }
	panic("not implemented")
}

func (es *EmailService) renderBatchedPost(notification *batchedNotification, sender *model.User, siteURL string, displayNameFormat string, translateFunc i18n.TranslateFunc, userLocale string, emailNotificationContentsType string) (string, error) {
	panic("not implemented")
}
