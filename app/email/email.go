package email

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/mail"
	"github.com/sitename/sitename/modules/templates"
)

func (es *Service) SendChangeUsernameEmail(newUsername, email, locale, siteURL string) error {
	T := i18n.GetUserTranslations(locale)

	subject := T(
		"api.templates.username_change_subject",
		map[string]interface{}{
			"Sitename": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.username_change_body.title")
	data.Props["Info"] = T(
		"api.templates.username_change_body.info",
		map[string]interface{}{
			"TeamDisplayName": es.config().ServiceSettings.SiteURL,
			"NewUsername":     newUsername,
		},
	)
	data.Props["Warning"] = T("api.templates.email_warning")

	body, err := es.templatesContainer.RenderToString("email_change_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) SendEmailChangeVerifyEmail(newUserEmail, locale, siteURL, token string) error {
	T := i18n.GetUserTranslations(locale)

	link := fmt.Sprintf("%s/do_verify_email?token=%s&email=%s", siteURL, token, url.QueryEscape(newUserEmail))

	subject := T(
		"api.templates.email_change_verify_subject",
		map[string]interface{}{
			"SiteName":        es.config().ServiceSettings.SiteURL,
			"TeamDisplayName": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.email_change_verify_body.title")
	data.Props["Info"] = T(
		"api.templates.email_change_verify_body.info",
		map[string]interface{}{
			"TeamDisplayName": es.config().ServiceSettings.SiteURL,
		},
	)
	data.Props["VerifyUrl"] = link
	data.Props["VerifyButton"] = T("api.templates.email_change_verify_body.button")

	body, err := es.templatesContainer.RenderToString("email_change_verify_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(newUserEmail, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) SendEmailChangeEmail(oldEmail, newEmail, locale, siteURL string) error {
	T := i18n.GetUserTranslations(locale)

	subject := T(
		"api.templates.email_change_subject",
		map[string]interface{}{
			"SiteName":        es.config().ServiceSettings.SiteURL,
			"TeamDisplayName": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.email_change_body.title")
	data.Props["Info"] = T(
		"api.templates.email_change_body.info",
		map[string]interface{}{
			"TeamDisplayName": es.config().ServiceSettings.SiteURL,
			"NewEmail":        newEmail,
		},
	)
	data.Props["Warning"] = T("api.templates.email_warning")

	body, err := es.templatesContainer.RenderToString("email_change_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(oldEmail, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) SendVerifyEmail(userEmail, locale, siteURL, token, redirect string) error {
	T := i18n.GetUserTranslations(locale)

	link := fmt.Sprintf("%s/do_verify_email?token=%s&email=%s", siteURL, token, url.QueryEscape(userEmail))
	if redirect != "" {
		link += fmt.Sprintf("&redirect_to=%s", redirect)
	}

	serverURL := condenseSiteURL(siteURL)

	subject := T("api.templates.verify_subject",
		map[string]interface{}{
			"SiteName": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.verify_body.title")
	data.Props["SubTitle1"] = T("api.templates.verify_body.subTitle1")
	data.Props["ServerURL"] = T("api.templates.verify_body.serverURL", map[string]interface{}{"ServerURL": serverURL})
	data.Props["SubTitle2"] = T("api.templates.verify_body.subTitle2")
	data.Props["ButtonURL"] = link
	data.Props["Button"] = T("api.templates.verify_body.button")
	data.Props["Info"] = T("api.templates.verify_body.info")
	data.Props["Info1"] = T("api.templates.verify_body.info1")
	data.Props["QuestionTitle"] = T("api.templates.questions_footer.title")
	data.Props["QuestionInfo"] = T("api.templates.questions_footer.info")

	body, err := es.templatesContainer.RenderToString("verify_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(userEmail, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) SendSignInChangeEmail(email, method, locale, siteURL string) error {
	T := i18n.GetUserTranslations(locale)

	subject := T(
		"api.templates.signin_change_email.subject",
		map[string]interface{}{
			"SiteName": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.signin_change_email.body.title")
	data.Props["Info"] = T(
		"api.templates.signin_change_email.body.info",
		map[string]interface{}{
			"SiteName": es.config().ServiceSettings.SiteURL,
			"Method":   method,
		},
	)
	data.Props["Warning"] = T("api.templates.email_warning")

	body, err := es.templatesContainer.RenderToString("signin_change_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) SendWelcomeEmail(userID string, email string, verified bool, disableWelcomeEmail bool, locale, siteURL, redirect string) error {
	if disableWelcomeEmail {
		return nil
	}
	if !*es.config().EmailSettings.SendEmailNotifications && !*es.config().EmailSettings.RequireEmailVerification {
		return errors.New("send email notifications and require email verification is disabled in the system console")
	}

	T := i18n.GetUserTranslations(locale)

	serverURL := condenseSiteURL(siteURL)

	subject := T(
		"api.templates.welcome_subject",
		map[string]interface{}{
			"SiteName":  es.config().ServiceSettings.SiteURL,
			"ServerURL": serverURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.welcome_body.title")
	data.Props["SubTitle1"] = T("api.templates.welcome_body.subTitle1")
	data.Props["ServerURL"] = T("api.templates.welcome_body.serverURL", map[string]interface{}{"ServerURL": serverURL})
	data.Props["SubTitle2"] = T("api.templates.welcome_body.subTitle2")
	data.Props["Button"] = T("api.templates.welcome_body.button")
	data.Props["Info"] = T("api.templates.welcome_body.info")
	data.Props["Info1"] = T("api.templates.welcome_body.info1")
	data.Props["SiteURL"] = siteURL

	if *es.config().NativeAppSettings.AppDownloadLink != "" {
		data.Props["AppDownloadTitle"] = T("api.templates.welcome_body.app_download_title")
		data.Props["AppDownloadInfo"] = T("api.templates.welcome_body.app_download_info")
		data.Props["AppDownloadButton"] = T("api.templates.welcome_body.app_download_button")
		data.Props["AppDownloadLink"] = *es.config().NativeAppSettings.AppDownloadLink
	}

	if !verified && *es.config().EmailSettings.RequireEmailVerification {
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

	body, err := es.templatesContainer.RenderToString("welcome_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) SendUserAccessTokenAddedEmail(email, locale, siteURL string) error {
	T := i18n.GetUserTranslations(locale)

	subject := T(
		"api.templates.user_access_token_subject",
		map[string]interface{}{
			"SiteName": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.user_access_token_body.title")
	data.Props["Info"] = T(
		"api.templates.user_access_token_body.info",
		map[string]interface{}{
			"SiteName": es.config().ServiceSettings.SiteURL,
			"SiteURL":  siteURL,
		},
	)
	data.Props["Warning"] = T("api.templates.email_warning")

	body, err := es.templatesContainer.RenderToString("password_change_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) SendPasswordChangeEmail(email, method, locale, siteURL string) error {
	T := i18n.GetUserTranslations(locale)

	subject := T(
		"api.templates.password_change_subject",
		map[string]interface{}{
			"SiteName":        es.config().ServiceSettings.SiteURL,
			"TeamDisplayName": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.password_change_body.title")
	data.Props["Info"] = T(
		"api.templates.password_change_body.info",
		map[string]interface{}{
			"TeamDisplayName": es.config().ServiceSettings.SiteURL,
			"TeamURL":         siteURL,
			"Method":          method,
		},
	)
	data.Props["Warning"] = T("api.templates.email_warning")

	body, err := es.templatesContainer.RenderToString("password_change_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) SendPasswordResetEmail(email string, token *model.Token, locale, siteURL string) (bool, error) {
	T := i18n.GetUserTranslations(locale)

	link := fmt.Sprintf("%s/reset_password_complete?token=%s", siteURL, url.QueryEscape(token.Token))

	subject := T(
		"api.templates.reset_subject",
		map[string]interface{}{
			"SiteName": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL
	data.Props["Title"] = T("api.templates.reset_body.title")
	data.Props["SubTitle"] = T("api.templates.reset_body.subTitle")
	data.Props["Info"] = T("api.templates.reset_body.info")
	data.Props["ButtonURL"] = link
	data.Props["Button"] = T("api.templates.reset_body.button")
	data.Props["QuestionTitle"] = T("api.templates.questions_footer.title")
	data.Props["QuestionInfo"] = T("api.templates.questions_footer.info")

	body, err := es.templatesContainer.RenderToString("reset_body", data)
	if err != nil {
		return false, err
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return false, err
	}

	return true, nil
}

func (es *Service) SendMfaChangeEmail(email string, activated bool, locale, siteURL string) error {
	T := i18n.GetUserTranslations(locale)

	subject := T(
		"api.templates.mfa_change_subject",
		map[string]interface{}{
			"SiteName": es.config().ServiceSettings.SiteURL,
		},
	)

	data := es.NewEmailTemplateData(locale)
	data.Props["SiteURL"] = siteURL

	if activated {
		data.Props["Info"] = T("api.templates.mfa_activated_body.info", map[string]interface{}{"SiteURL": siteURL})
		data.Props["Title"] = T("api.templates.mfa_activated_body.title")
	} else {
		data.Props["Info"] = T("api.templates.mfa_deactivated_body.info", map[string]interface{}{"SiteURL": siteURL})
		data.Props["Title"] = T("api.templates.mfa_deactivated_body.title")
	}
	data.Props["Warning"] = T("api.templates.email_warning")

	body, err := es.templatesContainer.RenderToString("mfa_change_body", data)
	if err != nil {
		return err
	}

	if err := es.sendMail(email, subject, body); err != nil {
		return err
	}

	return nil
}

func (es *Service) CreateVerifyEmailToken(userID string, newEmail string) (*model.Token, error) {
	tokenExtra := struct {
		UserId string
		Email  string
	}{
		userID,
		newEmail,
	}

	jsonData, err := json.Marshal(tokenExtra)
	if err != nil {
		return nil, errors.Wrap(CreateEmailTokenError, err.Error())
	}

	token := model_helper.NewToken(model_helper.TokenTypeVerifyEmail, string(jsonData))

	if err := es.InvalidateVerifyEmailTokensForUser(userID); err != nil {
		return nil, err
	}

	savedToken, err := es.store.Token().Save(*token)
	if err != nil {
		return nil, err
	}

	return savedToken, nil
}

func (es *Service) InvalidateVerifyEmailTokensForUser(userID string) *model_helper.AppError {
	tokens, err := es.store.Token().GetAllTokensByType(model_helper.TokenTypeVerifyEmail)
	if err != nil {
		return model_helper.NewAppError("InvalidateVerifyEmailTokensForUser", "api.user.invalidate_verify_email_tokens.error", nil, err.Error(), http.StatusInternalServerError)
	}

	var appErr *model_helper.AppError
	for _, token := range tokens {
		tokenExtra := struct {
			UserId string
			Email  string
		}{}

		if err := json.Unmarshal([]byte(token.Extra), &tokenExtra); err != nil {
			appErr = model_helper.NewAppError("InvalidateVerifyEmailTokensForUser", "api.user.invalidate_verify_email_tokens_parse.error", nil, err.Error(), http.StatusInternalServerError)
			continue
		}

		if tokenExtra.UserId != userID {
			continue
		}

		if err := es.store.Token().Delete(token.Token); err != nil {
			appErr = model_helper.NewAppError("InvalidateVerifyEmailTokensForUser", "api.user.invalidate_verify_email_tokens_delete.error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return appErr
}

func (es *Service) NewEmailTemplateData(locale string) templates.Data {
	var localT i18n.TranslateFunc
	if locale != "" {
		localT = i18n.GetUserTranslations(locale)
	} else {
		localT = i18n.T
	}
	organization := ""

	if feedBackOrganization := *es.config().EmailSettings.FeedbackOrganization; feedBackOrganization != "" {
		organization = localT("api.templates.email_organization") + feedBackOrganization
	}

	return templates.Data{
		Props: map[string]interface{}{
			"EmailInfo1": localT("api.templates.email_info1"),
			"EmailInfo2": localT("api.templates.email_info2"),
			"EmailInfo3": localT(
				"api.templates.email_info3",
				map[string]interface{}{
					"SiteName": es.config().ServiceSettings.SiteURL,
				},
			),
			"SupportEmail": *es.config().SupportSettings.SupportEmail,
			"Footer":       localT("api.templates.email_footer"),
			"FooterV2":     localT("api.templates.email_footer_v2"),
			"Organization": organization,
		},
		HTML: map[string]template.HTML{},
	}
}

func (es *Service) SendNotificationMail(to, subject, htmlBody string) error {
	if !*es.config().EmailSettings.SendEmailNotifications {
		return nil
	}
	return es.sendMail(to, subject, htmlBody)
}

func (es *Service) sendMail(to, subject, htmlBody string) error {
	return es.sendMailWithCC(to, subject, htmlBody, "")
}

func (es *Service) sendMailWithCC(to, subject, htmlBody string, ccMail string) error {
	mailConfig := es.mailServiceConfig()

	return mail.SendMailUsingConfig(to, subject, htmlBody, mailConfig, true, ccMail)
}
