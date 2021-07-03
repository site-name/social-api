package app

import (
	"fmt"
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
	TokenTypePasswordRecovery = "password_recovery" // token type for password recorver
	TokenTypeVerifyEmail      = "verify_email"      // type for creating user signup verification token
	TokenTypeGuestInvitation  = "guest_invitation"  // type for creating invite token
	TokenTypeCWSAccess        = "cws_access_token"
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

// server's email service
type EmailService struct {
	srv                     *Server
	PerHourEmailRateLimiter *throttled.GCRARateLimiter
	PerDayEmailRateLimiter  *throttled.GCRARateLimiter
	EmailBatching           *EmailBatchingJob
}

// NewEmailService create new email service and returns it
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
func (es *EmailService) SendChangeUsernameEmail(newUsername, email, locale, siteURL string) *model.AppError {
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

func (es *EmailService) SendEmailChangeVerifyEmail(newUserEmail, locale, siteURL, token string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendEmailChangeEmail(oldEmail, newEmail, locale, siteURL string) *model.AppError {
	T := i18n.GetUserTranslations(locale)

	subject := T("api.templates.email_change_subject",
		map[string]interface{}{
			"SiteName": model.TEAM_SETTINGS_DEFAULT_SITE_NAME,
		})

	data := es.newEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.email_change_body.title")
	data.Props["Info"] = T("api.templates.email_change_body.info",
		map[string]interface{}{
			"NewEmail": newEmail,
		})
	data.Props["Warning"] = T("api.templates.email_warning")

	body, err := es.srv.TemplatesContainer().RenderToString("email_change_body", data)
	if err != nil {
		return model.NewAppError("sendEmailChangeEmail", "api.user.send_email_change_email_and_forget.error", nil, err.Error(), http.StatusInternalServerError)
	}

	if err := es.sendMail(oldEmail, subject, body); err != nil {
		return model.NewAppError("sendEmailChangeEmail", "api.user.send_email_change_email_and_forget.error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (es *EmailService) SendVerifyEmail(userEmail, locale, siteURL, token, redirect string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendSignInChangeEmail(email, method, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

// after an user has signed up, send them an welcome email to newly registered user.
func (es *EmailService) SendWelcomeEmail(userID string, email string, verified bool, disableWelcomeEmail bool, locale, siteURL, redirect string) *model.AppError {
	if disableWelcomeEmail {
		return nil
	}
	if !*es.srv.Config().EmailSettings.SendEmailNotifications && !*es.srv.Config().EmailSettings.RequireEmailVerification {
		return model.NewAppError("SendWelcomeEmail", "api.user.send_welcome_email_and_forget.failed>error", nil, "Send Email Notifications and Require Email Verification is disabled", http.StatusInternalServerError)
	}

	T := i18n.GetUserTranslations(locale)

	serverURL := condenseSiteURL(siteURL)

	subject := T("api.templates.welcopme_subject", map[string]interface{}{
		"SiteName":  model.TEAM_SETTINGS_DEFAULT_SITE_NAME,
		"ServerURL": serverURL,
	})

	data := es.newEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.welcome_body.title")
	data.Props["SubTitle1"] = T("api.templates.welcome_body.subTitle1")
	data.Props["ServerURL"] = T("api.templates.welcome_body.serverURL", map[string]interface{}{"ServerURL": serverURL})
	data.Props["SubTitle2"] = T("api.templates.welcome_body.subTitle2")
	data.Props["Button"] = T("api.templates.welcome_body.button")
	data.Props["Info"] = T("api.templates.welcome_body.info")
	data.Props["Info1"] = T("api.templates.welcome_body.info1")

	if *es.srv.Config().NativeAppSettings.AppDownloadLink != "" {
		data.Props["AppDownloadTitle"] = T("api.templates.welcome_body.app_download_title")
		data.Props["AppDownloadInfo"] = T("api.templates.welcome_body.app_download_info")
		data.Props["AppDownloadButton"] = T("api.templates.welcome_body.app_download_button")
		data.Props["AppDownloadLink"] = *es.srv.Config().NativeAppSettings.AppDownloadLink
	}

	if !verified && *es.srv.Config().EmailSettings.RequireEmailVerification {
		token, err := es.CreateVerifyEmailToken(userID, email)
		if err != nil {
			return err
		}
		link := fmt.Sprintf("%s/do_verify_email?token=%s&email=%s", siteURL, token.Token, url.QueryEscape(email))
		if redirect != "" {
			link += fmt.Sprintf("&redirect_to=%s", redirect)
		}
		data.Props["ButtonURL"] = link
	}

	body, err := es.srv.TemplatesContainer().RenderToString("welcome_body", data)
	if err != nil {
		return model.NewAppError("sendWelcomeEmail", "api.user.send_welcome_email_and_forget.failed.error", nil, err.Error(), http.StatusInternalServerError)
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return model.NewAppError("sendWelcomeEmail", "api.user.send_welcome_email_and_forget.failed.error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (es *EmailService) SendPasswordChangeEmail(email, method, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendUserAccessTokenAddedEmail(email, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendPasswordResetEmail(email string, token *model.Token, locale, siteURL string) (bool, *model.AppError) {
	T := i18n.GetUserTranslations(locale)

	link := fmt.Sprintf("%s/reset_password_complete?token=%s", siteURL, url.QueryEscape(token.Token))

	subject := T("api.templates.reset_subject", map[string]interface{}{
		"Sitename": model.TEAM_SETTINGS_DEFAULT_SITE_NAME,
	})

	data := es.newEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.reset_body.title")
	data.Props["SubTitle"] = T("api.templates.reset_body.subTitle")
	data.Props["Info"] = T("api.templates.reset_body.info")
	data.Props["ButtonURL"] = link
	data.Props["Button"] = T("api.templates.reset_body.button")
	data.Props["QuestionTitle"] = T("api.templates.questions_footer.title")
	data.Props["QuestionInfo"] = T("api.templates.questions_footer.info")

	body, err := es.srv.TemplatesContainer().RenderToString("reset_body", data)
	if err != nil {
		return false, model.NewAppError("SendPasswordReset", "api.user.send_password_reset.send.app_error", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return false, model.NewAppError("SendPasswordReset", "api.user.send_password_reset.send.app_error", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	return true, nil
}

func (es *EmailService) SendMfaChangeEmail(email string, activated bool, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendDeactivateAccountEmail(email string, locale, siteURL string) *model.AppError {
	panic("not implemented")

}

func (es *EmailService) SendNotificationMail(to, subject, htmlBody string) error {
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

// sendMail is common method for sending emails
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

func (es *EmailService) SendMailWithEmbeddedFiles(to, subject, htmlBody string, embeddedFiles map[string]io.Reader) error {
	panic("not implemented")

}

type tokenExtra struct {
	UserId string
	Email  string
}

// CreateVerifyEmailToken create verification token
func (es *EmailService) CreateVerifyEmailToken(userID string, newEmail string) (*model.Token, *model.AppError) {
	tokenExtra := tokenExtra{
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

func (es *EmailService) SendNoCardPaymentFailedEmail(email string, locale string, siteURL string) *model.AppError {
	panic("not implemented")

}
