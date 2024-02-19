package app

import (
	"net/http"
	"strings"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/mail"
	"github.com/sitename/sitename/modules/slog"
)

func (a *App) GetWarnMetricsStatus() (map[string]*model_helper.WarnMetricStatus, *model_helper.AppError) {
	systemDataList, nErr := a.Srv().Store.System().Get()
	if nErr != nil {
		return nil, model_helper.NewAppError("GetWarnMetricsStatus", "app.system.get.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}

	isE0Edition := model_helper.BuildEnterpriseReady == "true" // license == nil was already validated upstream

	result := map[string]*model_helper.WarnMetricStatus{}
	for key, value := range systemDataList {
		if strings.HasPrefix(key, model_helper.WarnMetricStatusStorePrefix) {
			if warnMetric, ok := model_helper.WarnMetricsTable[key]; ok {
				if !warnMetric.IsBotOnly && (value == model_helper.WarnMetricStatusRunonce || value == model_helper.WarnMetricStatusLimitReached) {
					result[key], _ = a.getWarnMetricStatusAndDisplayTextsForId(key, nil, isE0Edition)
				}
			}
		}
	}

	return result, nil
}

func (a *App) getWarnMetricStatusAndDisplayTextsForId(warnMetricId string, T i18n.TranslateFunc, isE0Edition bool) (*model_helper.WarnMetricStatus, *model_helper.WarnMetricDisplayTexts) {
	var warnMetricStatus *model_helper.WarnMetricStatus
	var warnMetricDisplayTexts = &model_helper.WarnMetricDisplayTexts{}

	if warnMetric, ok := model_helper.WarnMetricsTable[warnMetricId]; ok {
		warnMetricStatus = &model_helper.WarnMetricStatus{
			Id:    warnMetric.Id,
			Limit: warnMetric.Limit,
			Acked: false,
		}

		if T == nil {
			slog.Debug("No translation function")
			return warnMetricStatus, nil
		}

		warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.bot_response.notification_success.message")

		switch warnMetricId {
		case model_helper.SystemWarnMetricNumberOfTeams5:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.number_of_teams_5.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_teams_5.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.number_of_teams_5.start_trial_notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.number_of_teams_5.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_teams_5.notification_body")
			}
		case model_helper.SystemWarnMetricMfa:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.mfa.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.mfa.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.mfa.start_trial_notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.mfa.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.mfa.notification_body")
			}
		case model_helper.SystemWarnMetricEmailDomain:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.email_domain.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.email_domain.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.email_domain.start_trial_notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.email_domain.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.email_domain.notification_body")
			}
		case model_helper.SystemWarnMetricNumberOfChannels50:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.number_of_channels_50.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_channels_50.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.number_of_channels_50.start_trial.notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.number_of_channels_50.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_channels_50.notification_body")
			}
		case model_helper.SystemWarnMetricNumberOfActiveUsers100:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.number_of_active_users_100.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_active_users_100.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.number_of_active_users_100.start_trial.notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.number_of_active_users_100.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_active_users_100.notification_body")
			}
		case model_helper.SystemWarnMetricNumberOfActiveUsers200:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.number_of_active_users_200.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_active_users_200.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.number_of_active_users_200.start_trial.notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.number_of_active_users_200.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_active_users_200.notification_body")
			}
		case model_helper.SystemWarnMetricNumberOfActiveUsers300:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.number_of_active_users_300.start_trial.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_active_users_300.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.number_of_active_users_300.start_trial.notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.number_of_active_users_300.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_active_users_300.notification_body")
			}
		case model_helper.SystemWarnMetricNumberOfActiveUsers500:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.number_of_active_users_500.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_active_users_500.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.number_of_active_users_500.start_trial.notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.number_of_active_users_500.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_active_users_500.notification_body")
			}
		case model_helper.SystemWarnMetricNumberOfPosts2m:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.number_of_posts_2M.notification_title")
			if isE0Edition {
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_posts_2M.start_trial.notification_body")
				warnMetricDisplayTexts.BotSuccessMessage = T("api.server.warn_metric.number_of_posts_2M.start_trial.notification_success.message")
			} else {
				warnMetricDisplayTexts.EmailBody = T("api.server.warn_metric.number_of_posts_2M.contact_us.email_body")
				warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.number_of_posts_2M.notification_body")
			}
		case model_helper.SystemMetricSupportEmailNotConfigured:
			warnMetricDisplayTexts.BotTitle = T("api.server.warn_metric.support_email_not_configured.notification_title")
			warnMetricDisplayTexts.BotMessageBody = T("api.server.warn_metric.support_email_not_configured.start_trial.notification_body")
		default:
			slog.Debug("Invalid metric id", slog.String("id", warnMetricId))
			return nil, nil
		}

		return warnMetricStatus, warnMetricDisplayTexts
	}
	return nil, nil
}

//nolint:golint,unused,deadcode
// func (a *App) notifyAdminsOfWarnMetricStatus(c *request.Context, warnMetricId string, isE0Edition bool) *model_helper.AppError {
// 	// get warn metrics bot
// 	warnMetricsBot, err := a.GetWarnMetricsBot()
// 	if err != nil {
// 		return err
// 	}

// 	warnMetric, ok := model.WarnMetricsTable[warnMetricId]
// 	if !ok {
// 		return model_helper.NewAppError("NotifyAdminsOfWarnMetricStatus", "app.system.warn_metric.notification.invalid_metric.app_error", nil, "", http.StatusInternalServerError)
// 	}

// 	perPage := 25
// 	userOptions := &account.UserGetOptions{
// 		Page:     0,
// 		PerPage:  perPage,
// 		Role:     model.SYSTEM_ADMIN_ROLE_ID,
// 		Inactive: false,
// 	}

// 	// get sysadmins
// 	var sysAdmins []*account.User
// 	for {
// 		sysAdminsList, err := a.account.GetUsers(userOptions)
// 		if err != nil {
// 			return err
// 		}

// 		if len(sysAdminsList) == 0 {
// 			return model_helper.NewAppError("NotifyAdminsOfWarnMetricStatus", "app.system.warn_metric.notification.empty_admin_list.app_error", nil, "", http.StatusInternalServerError)
// 		}
// 		sysAdmins = append(sysAdmins, sysAdminsList...)

// 		if len(sysAdminsList) < perPage {
// 			slog.Debug("Number of system admins is less than page limit", slog.Int("count", len(sysAdminsList)))
// 			break
// 		}
// 	}

// 	for _, sysAdmin := range sysAdmins {
// 		T := i18n.GetUserTranslations(sysAdmin.Locale)
// 		warnMetricsBot.DisplayName = T("app.system.warn_metric.bot_displayname")
// 		warnMetricsBot.Description = T("app.system.warn_metric.bot_description")

// 		channel, appErr := a.GetOrCreateDirectChannel(c, warnMetricsBot.UserId, sysAdmin.Id)
// 		if appErr != nil {
// 			return appErr
// 		}

// 		warnMetricStatus, warnMetricDisplayTexts := a.getWarnMetricStatusAndDisplayTextsForId(warnMetricId, T, isE0Edition)
// 		if warnMetricStatus == nil {
// 			return model_helper.NewAppError("NotifyAdminsOfWarnMetricStatus", "app.system.warn_metric.notification.invalid_metric.app_error", nil, "", http.StatusInternalServerError)
// 		}

// 		botPost := &model.Post{
// 			UserId:    warnMetricsBot.UserId,
// 			ChannelId: channel.Id,
// 			Type:      model.PostTypeSystemWarnMetricStatus,
// 			Message:   "",
// 		}

// 		actionId := "contactUs"
// 		actionName := T("api.server.warn_metric.contact_us")
// 		postActionValue := T("api.server.warn_metric.contacting_us")
// 		postActionUrl := fmt.Sprintf("/warn_metrics/ack/%s", warnMetricId)

// 		if isE0Edition {
// 			actionId = "startTrial"
// 			actionName = T("api.server.warn_metric.start_trial")
// 			postActionValue = T("api.server.warn_metric.starting_trial")
// 			postActionUrl = fmt.Sprintf("/warn_metrics/trial-license-ack/%s", warnMetricId)
// 		}

// 		actions := []*model.PostAction{}
// 		actions = append(actions,
// 			&model.PostAction{
// 				Id:   actionId,
// 				Name: actionName,
// 				Type: model.PostActionTypeButton,
// 				Options: []*model.PostActionOptions{
// 					{
// 						Text:  "TrackEventId",
// 						Value: warnMetricId,
// 					},
// 					{
// 						Text:  "ActionExecutingMessage",
// 						Value: postActionValue,
// 					},
// 				},
// 				Integration: &model.PostActionIntegration{
// 					Context: model.StringInterface{
// 						"bot_user_id": warnMetricsBot.UserId,
// 						"force_ack":   false,
// 					},
// 					URL: postActionUrl,
// 				},
// 			},
// 		)

// 		attachments := []*model.SlackAttachment{{
// 			AuthorName: "",
// 			Title:      warnMetricDisplayTexts.BotTitle,
// 			Text:       warnMetricDisplayTexts.BotMessageBody,
// 		}}

// 		if !warnMetric.SkipAction {
// 			attachments[0].Actions = actions
// 		}

// 		model.ParseSlackAttachment(botPost, attachments)

// 		slog.Debug("Post admin advisory for metric", slog.String("warnMetricId", warnMetricId), slog.String("userid", botPost.UserId))
// 		if _, err := a.CreatePostAsUser(c, botPost, c.Session().Id, true); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (a *App) NotifyAndSetWarnMetricAck(warnMetricId string, sender model.User, forceAck bool, isBot bool) *model_helper.AppError {
	if warnMetric, ok := model_helper.WarnMetricsTable[warnMetricId]; ok {
		data, nErr := a.Srv().Store.System().GetByName(warnMetric.Id)
		if nErr == nil && data != nil && data.Value == model_helper.WarnMetricStatusAck {
			slog.Debug("This metric warning has already been acknowledged", slog.String("id", warnMetric.Id))
			return nil
		}

		if !forceAck {
			if *a.Config().EmailSettings.SMTPServer == "" {
				return model_helper.NewAppError("NotifyAndSetWarnMetricAck", "api.email.send_warn_metric_ack.missing_server.app_error", nil, i18n.T("api.context.invalid_param.app_error", map[string]interface{}{"Name": "SMTPServer"}), http.StatusInternalServerError)
			}
			T := i18n.GetUserTranslations(sender.Locale.String())
			data := a.Srv().EmailService.NewEmailTemplateData(sender.Locale.String())
			data.Props["ContactNameHeader"] = T("api.templates.warn_metric_ack.body.contact_name_header")
			data.Props["ContactNameValue"] = model_helper.UserGetFullName(sender)
			data.Props["ContactEmailHeader"] = T("api.templates.warn_metric_ack.body.contact_email_header")
			data.Props["ContactEmailValue"] = sender.Email

			//same definition as the active users count metric displayed in the SystemConsole Analytics section
			registeredUsersCount, cerr := a.Srv().Store.User().Count(model_helper.UserCountOptions{})
			if cerr != nil {
				slog.Warn("Error retrieving the number of registered users", slog.Err(cerr))
			} else {
				data.Props["RegisteredUsersHeader"] = T("api.templates.warn_metric_ack.body.registered_users_header")
				data.Props["RegisteredUsersValue"] = registeredUsersCount
			}
			data.Props["SiteURLHeader"] = T("api.templates.warn_metric_ack.body.site_url_header")
			data.Props["SiteURL"] = a.GetSiteURL()
			data.Props["TelemetryIdHeader"] = T("api.templates.warn_metric_ack.body.diagnostic_id_header")
			// data.Props["TelemetryIdValue"] = a.TelemetryId()
			data.Props["Footer"] = T("api.templates.warn_metric_ack.footer")

			warnMetricStatus, warnMetricDisplayTexts := a.getWarnMetricStatusAndDisplayTextsForId(warnMetricId, T, false)
			if warnMetricStatus == nil {
				return model_helper.NewAppError("NotifyAndSetWarnMetricAck", "api.email.send_warn_metric_ack.invalid_warn_metric.app_error", nil, "", http.StatusInternalServerError)
			}

			subject := T("api.templates.warn_metric_ack.subject")
			data.Props["Title"] = warnMetricDisplayTexts.EmailBody

			mailConfig := a.Srv().MailServiceConfig()

			body, err := a.Srv().TemplatesContainer().RenderToString("warn_metric_ack", data)
			if err != nil {
				return model_helper.NewAppError("NotifyAndSetWarnMetricAck", "api.email.send_warn_metric_ack.failure.app_error", map[string]interface{}{"Error": err.Error()}, "", http.StatusInternalServerError)
			}

			if err := mail.SendMailUsingConfig(model_helper.MM_SUPPORT_ADVISOR_ADDRESS, subject, body, mailConfig, false, sender.Email); err != nil {
				return model_helper.NewAppError("NotifyAndSetWarnMetricAck", "api.email.send_warn_metric_ack.failure.app_error", map[string]interface{}{"Error": err.Error()}, "", http.StatusInternalServerError)
			}
		}

		if err := a.setWarnMetricsStatusAndNotify(warnMetric.Id); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) setWarnMetricsStatusAndNotify(warnMetricId string) *model_helper.AppError {
	// Ack all metric warnings on the server
	if err := a.setWarnMetricsStatus(model_helper.WarnMetricStatusAck); err != nil {
		return err
	}

	// Inform client that this metric warning has been acked
	message := model_helper.NewWebSocketEvent(model_helper.WebsocketWarnMetricStatusRemoved, "", nil)
	message.Add("warnMetricId", warnMetricId)
	a.Publish(message)

	return nil
}

func (a *App) setWarnMetricsStatus(status string) *model_helper.AppError {
	slog.Debug("Set monitoring status for all warn metrics", slog.String("status", status))
	for _, warnMetric := range model_helper.WarnMetricsTable {
		if err := a.setWarnMetricsStatusForId(warnMetric.Id, status); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) setWarnMetricsStatusForId(warnMetricId string, status string) *model_helper.AppError {
	slog.Debug("Store status for warn metric", slog.String("warnMetricId", warnMetricId), slog.String("status", status))
	if err := a.Srv().Store.System().SaveOrUpdateWithWarnMetricHandling(model.System{
		Name:  warnMetricId,
		Value: status,
	}); err != nil {
		return model_helper.NewAppError("setWarnMetricsStatusForId", "app.system.warn_metric.store.app_error", map[string]interface{}{"WarnMetricName": warnMetricId}, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
