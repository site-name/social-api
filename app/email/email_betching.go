package email

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
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

func (es *Service) AddNotificationEmailToBatch(user *account.User, post *model.Post, team *model.Team) *model.AppError {
	if !*es.config().EmailSettings.EnableEmailBatching {
		return model.NewAppError("AddNotificationEmailToBatch", "api.email_batching.add_notification_email_to_batch.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if !es.EmailBatching.Add(user, post, team) {
		slog.Error("Email batching job's receiving channel was full. Please increase the EmailBatchingBufferSize.")
		return model.NewAppError("AddNotificationEmailToBatch", "api.email_batching.add_notification_email_to_batch.channel_full.app_error", nil, "", http.StatusInternalServerError)
	}

	return nil
}

type batchedNotification struct {
	userID string
	// post     *model.Post
	// teamName string
}

type EmailBatchingJob struct {
	config               func() *model.Config
	service              *Service
	newNotifications     chan *batchedNotification
	pendingNotifications map[string][]*batchedNotification
	task                 *model.ScheduledTask
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
	newTask := model.CreateRecurringTask(EmailBatchingTaskName, job.CheckPendingEmails, time.Duration(*job.config().EmailSettings.EmailBatchingInterval)*time.Second)

	job.taskMutex.Lock()
	oldTask := job.task
	job.task = newTask
	job.taskMutex.Unlock()

	if oldTask != nil {
		oldTask.Cancel()
	}
}

// func(job *EmailBatchingJob) Add(user *)
