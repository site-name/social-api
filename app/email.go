package app

import (
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/mail"
	"github.com/sitename/sitename/modules/templates"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

const (
	emailRateLimitingMemstoreSize = 65536
	emailRateLimitingPerHour      = 20
	emailRateLimitingMaxBurst     = 20
)

func condenseSiteURL(siteURL string) string {
	parsedSiteURL, _ := url.Parse(siteURL)
	if parsedSiteURL.Path == "" || parsedSiteURL.Path == "/" {
		return parsedSiteURL.Host
	}

	return path.Join(parsedSiteURL.Host, parsedSiteURL.Path)
}

type EmailService struct {
	srv                     *Server
	PerHourEmailRateLimiter *throttled.GCRARateLimiter
	PerDayEmailRateLimiter  *throttled.GCRARateLimiter
	EmailBatching           *EmailBatchingJob
}

func NewEmailService(srv *Server) (*EmailService, error) {
	service := &EmailService{srv: srv}
	if err := service.setUpRateLimiters(); err != nil {
		return nil, err
	}
	service.InitEmailBatching()
	return service, nil
}

func (es *EmailService) setUpRateLimiters() error {
	store, err := memstore.New(emailRateLimitingMemstoreSize)
	if err != nil {
		return errors.Wrap(err, "Unable to setup email rate limiting memstore.")
	}

	perHourQuota := throttled.RateQuota{
		MaxRate:  throttled.PerHour(emailRateLimitingPerHour),
		MaxBurst: emailRateLimitingMaxBurst,
	}

	perDayQuota := throttled.RateQuota{
		MaxRate:  throttled.PerDay(1),
		MaxBurst: 0,
	}

	perHourRateLimiter, err := throttled.NewGCRARateLimiter(store, perHourQuota)
	if err != nil || perHourRateLimiter == nil {
		return errors.Wrap(err, "Unable to setup email rate limiting GCRA rate limiter.")
	}
	perDayRateLimiter, err := throttled.NewGCRARateLimiter(store, perDayQuota)
	if err != nil || perDayRateLimiter == nil {
		return errors.Wrap(err, "Unable to setup per day email rate limiting GCRA rate limiter.")
	}

	es.PerHourEmailRateLimiter = perHourRateLimiter
	es.PerDayEmailRateLimiter = perDayRateLimiter
	return nil
}

// sendChangeUsernameEmail send email to let user know username changes
func (es *EmailService) sendChangeUsernameEmail(newUsername, email, locale, siteURL string) *model.AppError {
	T := i18n.GetUserTranslations(locale)

	subject := T(
		"api.templates.username_change_subject",
		map[string]interface{}{
			"SiteName": model.TEAM_SETTINGS_DEFAULT_SITE_NAME,
		},
	)

	data := es.newEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.username_change_body.title")
	data.Props["Info"] = T("api.templates.username_change_body.info",
		map[string]interface{}{"TeamDisplayName": model.TEAM_SETTINGS_DEFAULT_SITE_NAME, "NewUsername": newUsername})
	data.Props["Warning"] = T("api.templates.email_warning")

	body, err := es.srv.TemplatesContainer().RenderToString("email_change_body", data)
	if err != nil {
		return model.NewAppError("sendChangeUsernameEmail", "api.user.send_email_change_username_and_forget.error", nil, err.Error(), http.StatusInternalServerError)
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return model.NewAppError("sendChangeUsernameEmail", "api.user.send_email_change_username_and_forget.error", nil, err.Error(), http.StatusInternalServerError)

	}

	return nil
}

func (es *EmailService) sendEmailChangeVerifyEmail(newUserEmail, locale, siteURL, token string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) sendEmailChangeEmail(oldEmail, newEmail, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) sendVerifyEmail(userEmail, locale, siteURL, token, redirect string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendSignInChangeEmail(email, method, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) sendWelcomeEmail(userID string, email string, verified bool, disableWelcomeEmail bool, locale, siteURL, redirect string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) sendPasswordChangeEmail(email, method, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) sendUserAccessTokenAddedEmail(email, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendPasswordResetEmail(email string, token *model.Token, locale, siteURL string) (bool, *model.AppError) {
	panic("not implemented")

}

func (es *EmailService) sendMfaChangeEmail(email string, activated bool, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendDeactivateAccountEmail(email string, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) sendNotificationMail(to, subject, htmlBody string) error {
	panic("not implemented")

}

func (es *EmailService) newEmailTemplateData(locale string) templates.Data {
	var localT i18n.TranslateFunc
	if locale != "" {
		localT = i18n.GetUserTranslations(locale)
	} else {
		localT = i18n.T
	}
	organization := ""

	emailSettings := es.srv.Config().EmailSettings
	if *emailSettings.FeedbackOrganization != "" {
		organization = localT("api.templates.email_organization") + *emailSettings.FeedbackOrganization
	}

	return templates.Data{
		Props: map[string]interface{}{
			"EmailInfo1":   localT("api.templates.email_info1"),
			"EmailInfo2":   localT("api.templates.email_info2"),
			"EmailInfo3":   localT("api.templates.email_info3", map[string]interface{}{"SiteName": model.TEAM_SETTINGS_DEFAULT_SITE_NAME}),
			"SupportEmail": *es.srv.Config().SupportSettings.SupportEmail,
			"Footer":       localT("api.templates.email_footer"),
			"FooterV2":     localT("api.templates.email_footer_v2"),
			"Organization": organization,
		},
		HTML: map[string]template.HTML{},
	}
}

func (es *EmailService) sendMail(to, subject, htmlBody string) error {
	return es.sendMailWithCC(to, subject, htmlBody, "")
}

func (es *EmailService) sendMailWithCC(to, subject, htmlBody string, ccMail string) error {
	mailConfig := es.srv.MailServiceConfig()

	return mail.SendMailUsingConfig(to, subject, htmlBody, mailConfig,
		true,
		ccMail,
	)
}

func (es *EmailService) sendMailWithEmbeddedFiles(to, subject, htmlBody string, embeddedFiles map[string]io.Reader) error {
	panic("not implemented")

}

func (es *EmailService) CreateVerifyEmailToken(userID string, newEmail string) (*model.Token, *model.AppError) {
	tokenExtra := struct {
		UserId string
		Email  string
	}{
		userID,
		newEmail,
	}
	jsonData, err := json.JSON.Marshal(tokenExtra)

	if err != nil {
		return nil, model.NewAppError("CreateVerifyEmailToken", "api.user.create_email_token.error", nil, "", http.StatusInternalServerError)
	}

	token := model.NewToken(TokenTypeVerifyEmail, string(jsonData))

	if err = es.srv.Store.Token().Save(token); err != nil {
		var appErr *model.AppError
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		default:
			return nil, model.NewAppError("CreateVerifyEmailToken", "app.recover.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return token, nil
}

func (es *EmailService) SendOverUserLimitWarningEmail(email string, locale string, siteURL string) (bool, *model.AppError) {
	panic("not implemented")
}

func (es *EmailService) SendPaymentFailedEmail(email string, locale string, failedPayment *model.FailedPayment, siteURL string) (bool, *model.AppError) {
	panic("not implemented")

}

func (es *EmailService) SendNoCardPaymentFailedEmail(email string, locale string, siteURL string) *model.AppError {
	panic("not implemented")

}
