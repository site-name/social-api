package app

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"runtime"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/mail"
	"github.com/sitename/sitename/modules/slog"
)

const (
	PropSecurityURL      = "https://securityupdatecheck.mattermost.com"
	SecurityUpdatePeriod = 86400000 // 24 hours in milliseconds.

	PropSecurityID              = "id"
	PropSecurityBuild           = "b"
	PropSecurityEnterpriseReady = "be"
	PropSecurityDatabase        = "db"
	PropSecurityOS              = "os"
	PropSecurityUserCount       = "uc"
	PropSecurityTeamCount       = "tc"
	PropSecurityActiveUserCount = "auc"
	PropSecurityUnitTests       = "ut"
)

// DoSecurityUpdateCheck
func (s *Server) DoSecurityUpdateCheck() {
	if !*s.Config().ServiceSettings.EnableSecurityFixAlert {
		return
	}

	props, err := s.Store.System().Get()
	if err != nil {
		return
	}

	lastSecurityTime, _ := strconv.ParseInt(props[model.SystemLastSecurityTime], 10, 0)
	currentTime := model.GetMillis()

	if (currentTime - lastSecurityTime) > SecurityUpdatePeriod {
		slog.Debug("Checking for security update from Mattermost")

		v := url.Values{}

		// v.Set(PropSecurityID, s.TelemetryId())
		v.Set(PropSecurityBuild, model.CurrentVersion+"."+model.BuildNumber)
		v.Set(PropSecurityEnterpriseReady, model.BuildEnterpriseReady)
		v.Set(PropSecurityDatabase, *s.Config().SqlSettings.DriverName)
		v.Set(PropSecurityOS, runtime.GOOS)

		if props[model.SystemRanUnitTests] != "" {
			v.Set(PropSecurityUnitTests, "1")
		} else {
			v.Set(PropSecurityUnitTests, "0")
		}

		systemSecurityLastTime := &model.System{Name: model.SystemLastSecurityTime, Value: strconv.FormatInt(currentTime, 10)}
		if lastSecurityTime == 0 {
			s.Store.System().Save(systemSecurityLastTime)
		} else {
			s.Store.System().Update(systemSecurityLastTime)
		}

		if count, err := s.Store.User().Count(model.UserCountOptions{IncludeDeleted: true}); err == nil {
			v.Set(PropSecurityUserCount, strconv.FormatInt(count, 10))
		}

		if ucr, err := s.Store.Status().GetTotalActiveUsersCount(); err == nil {
			v.Set(PropSecurityActiveUserCount, strconv.FormatInt(ucr, 10))
		}

		res, err := http.Get(PropSecurityURL + "/security?" + v.Encode())
		if err != nil {
			slog.Error("Failed to get security update information from Mattermost.")
			return
		}

		defer res.Body.Close()

		var bulletins model.SecurityBulletins
		if jsonErr := json.NewDecoder(res.Body).Decode(&bulletins); jsonErr != nil {
			slog.Error("failed to decode JSON", slog.Err(jsonErr))
			return
		}

		for _, bulletin := range bulletins {
			if bulletin.AppliesToVersion == model.CurrentVersion {
				if props["SecurityBulletin_"+bulletin.Id] == "" {
					users, userErr := s.Store.User().GetSystemAdminProfiles()
					if userErr != nil {
						slog.Error("Failed to get system admins for security update information from Mattermost.")
						return
					}

					resBody, err := http.Get(PropSecurityURL + "/bulletins/" + bulletin.Id)
					if err != nil {
						slog.Error("Failed to get security bulletin details")
						return
					}

					body, err := io.ReadAll(resBody.Body)
					resBody.Body.Close()
					if err != nil || resBody.StatusCode != 200 {
						slog.Error("Failed to read security bulletin details")
						return
					}

					for _, user := range users {
						slog.Info("Sending security bulletin", slog.String("bulletin_id", bulletin.Id), slog.String("user_email", user.Email))
						mailConfig := s.MailServiceConfig()
						mail.SendMailUsingConfig(user.Email, i18n.T("mattermost.bulletin.subject"), string(body), mailConfig, true, "")
					}

					bulletinSeen := &model.System{Name: "SecurityBulletin_" + bulletin.Id, Value: bulletin.Id}
					s.Store.System().Save(bulletinSeen)
				}
			}
		}
	}
}
