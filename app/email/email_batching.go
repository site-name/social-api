package email

import (
	"html/template"
	"sync"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/slog"
)

const (
	EmailBatchingTaskName = "Email Batching"
)

type postData struct {
	SenderName  string
	ChannelName string
	Message     template.HTML
	MessageURL  string
	SenderPhoto string
	PostPhoto   string
	Time        string
}

func (es *Service) InitEmailBatching() {
	if *es.config().EmailSettings.EnableEmailBatching {
		if es.EmailBatching == nil {
			es.EmailBatching = NewEmailBatchingJob(es, *es.config().EmailSettings.EmailBatchingBufferSize)
		}

		// note that we don't support changing EmailBatchingBufferSize without restarting the server

		es.EmailBatching.Start()
	}
}

// If the name is longer than i characters, replace remaning characters with ...
func truncateUserNames(name string, i int) string {
	runes := []rune(name)
	if len(runes) > i {
		newString := string(runes[:i])
		return newString + "..."
	}

	return name
}

// func (es *Service) AddNotificationEmailToBatch(user *account.User, post *model.Post, team *model.Team) *model_helper.AppError {
// 	if !*es.config().EmailSettings.EnableEmailBatching {
// 		return model_helper.NewAppError("AddNotificationEmailToBatch", "api.email_batching.add_notification_email_to_batch.disabled.app_error", nil, "", http.StatusNotImplemented)
// 	}

// 	if !es.EmailBatching.Add(user, post, team) {
// 		slog.Error("Email batching job's receiving channel was full. Please increase the EmailBatchingBufferSize.")
// 		return model_helper.NewAppError("AddNotificationEmailToBatch", "api.email_batching.add_notification_email_to_batch.channel_full.app_error", nil, "", http.StatusInternalServerError)
// 	}

// 	return nil
// }

type batchedNotification struct {
	userID string
	// post     *model.Post
	// teamName string
}

type EmailBatchingJob struct {
	config               func() *model_helper.Config
	service              *Service
	newNotifications     chan *batchedNotification
	pendingNotifications map[string][]*batchedNotification
	task                 *model_helper.ScheduledTask
	taskMutex            sync.Mutex
}

func NewEmailBatchingJob(es *Service, bufferSize int) *EmailBatchingJob {
	return &EmailBatchingJob{
		config:               es.config,
		service:              es,
		newNotifications:     make(chan *batchedNotification, bufferSize),
		pendingNotifications: make(map[string][]*batchedNotification),
	}
}

func (job *EmailBatchingJob) Start() {
	slog.Debug("Email batching job starting. Checking for pending emails periodically.", slog.Int("interval_in_seconds", *job.config().EmailSettings.EmailBatchingInterval))
	newTask := model_helper.CreateRecurringTask(EmailBatchingTaskName, job.CheckPendingEmails, time.Duration(*job.config().EmailSettings.EmailBatchingInterval)*time.Second)

	job.taskMutex.Lock()
	oldTask := job.task
	job.task = newTask
	job.taskMutex.Unlock()

	if oldTask != nil {
		oldTask.Cancel()
	}
}

func (job *EmailBatchingJob) Add(user *model.User) bool {
	notification := &batchedNotification{
		userID: user.ID,
	}

	select {
	case job.newNotifications <- notification:
		return true

	default:
		// return false if we couldn't queue the email notification so that we can send an immediate email
		return false
	}
}

func (job *EmailBatchingJob) CheckPendingEmails() {
	// job.handleNewNotitifcations()

	// it's a bit weird to pass the send email function through here, but it makes it so that we can test
	// without actually sending emails
	// job.checkPendingNotifications(time.Now(), job.service)
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
	panic("not implt")
}
