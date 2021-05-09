package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sitename/sitename/model"
)

// GenerateClientConfig renders the given configuration for a client.
func GenerateClientConfig(c *model.Config) map[string]string {
	props := GenerateLimitedClientConfig(c)

	props["SiteURL"] = strings.TrimRight(*c.ServiceSettings.SiteURL, "/")
	props["EnableBotAccountCreation"] = strconv.FormatBool(*c.ServiceSettings.EnableBotAccountCreation)
	props["EnableOAuthServiceProvider"] = strconv.FormatBool(*c.ServiceSettings.EnableOAuthServiceProvider)
	props["GoogleDeveloperKey"] = *c.ServiceSettings.GoogleDeveloperKey
	props["EnableIncomingWebhooks"] = strconv.FormatBool(*c.ServiceSettings.EnableIncomingWebhooks)
	props["EnableOutgoingWebhooks"] = strconv.FormatBool(*c.ServiceSettings.EnableOutgoingWebhooks)
	props["EnableCommands"] = strconv.FormatBool(*c.ServiceSettings.EnableCommands)
	props["EnablePostUsernameOverride"] = strconv.FormatBool(*c.ServiceSettings.EnablePostUsernameOverride)
	props["EnablePostIconOverride"] = strconv.FormatBool(*c.ServiceSettings.EnablePostIconOverride)
	props["EnableUserAccessTokens"] = strconv.FormatBool(*c.ServiceSettings.EnableUserAccessTokens)
	props["EnableLinkPreviews"] = strconv.FormatBool(*c.ServiceSettings.EnableLinkPreviews)
	props["EnableTesting"] = strconv.FormatBool(*c.ServiceSettings.EnableTesting)
	props["EnableDeveloper"] = strconv.FormatBool(*c.ServiceSettings.EnableDeveloper)
	props["PostEditTimeLimit"] = fmt.Sprintf("%v", *c.ServiceSettings.PostEditTimeLimit)
	props["MinimumHashtagLength"] = fmt.Sprintf("%v", *c.ServiceSettings.MinimumHashtagLength)
	props["CloseUnusedDirectMessages"] = strconv.FormatBool(*c.ServiceSettings.CloseUnusedDirectMessages)
	props["EnablePreviewFeatures"] = strconv.FormatBool(*c.ServiceSettings.EnablePreviewFeatures)
	props["EnableTutorial"] = strconv.FormatBool(*c.ServiceSettings.EnableTutorial)
	props["ExperimentalEnableDefaultChannelLeaveJoinMessages"] = strconv.FormatBool(*c.ServiceSettings.ExperimentalEnableDefaultChannelLeaveJoinMessages)
	props["ExperimentalGroupUnreadChannels"] = *c.ServiceSettings.ExperimentalGroupUnreadChannels
	props["EnableSVGs"] = strconv.FormatBool(*c.ServiceSettings.EnableSVGs)
	props["EnableMarketplace"] = strconv.FormatBool(*c.PluginSettings.EnableMarketplace)
	props["EnableLatex"] = strconv.FormatBool(*c.ServiceSettings.EnableLatex)
	props["ExtendSessionLengthWithActivity"] = strconv.FormatBool(*c.ServiceSettings.ExtendSessionLengthWithActivity)
	props["ManagedResourcePaths"] = *c.ServiceSettings.ManagedResourcePaths

	// This setting is only temporary, so keep using the old setting name for the mobile and web apps
	props["ExperimentalEnablePostMetadata"] = "true"
	props["ExperimentalEnableClickToReply"] = strconv.FormatBool(*c.ExperimentalSettings.EnableClickToReply)

	props["ExperimentalCloudUserLimit"] = strconv.FormatInt(*c.ExperimentalSettings.CloudUserLimit, 10)
	props["ExperimentalCloudBilling"] = strconv.FormatBool(*c.ExperimentalSettings.CloudBilling)
	if *c.ServiceSettings.ExperimentalChannelOrganization || *c.ServiceSettings.ExperimentalGroupUnreadChannels != model.GROUP_UNREAD_CHANNELS_DISABLED {
		props["ExperimentalChannelOrganization"] = strconv.FormatBool(true)
	} else {
		props["ExperimentalChannelOrganization"] = strconv.FormatBool(false)
	}

	props["ExperimentalTimezone"] = strconv.FormatBool(*c.DisplaySettings.ExperimentalTimezone)

	props["SendEmailNotifications"] = strconv.FormatBool(*c.EmailSettings.SendEmailNotifications)
	props["SendPushNotifications"] = strconv.FormatBool(*c.EmailSettings.SendPushNotifications)
	props["RequireEmailVerification"] = strconv.FormatBool(*c.EmailSettings.RequireEmailVerification)
	props["EnableEmailBatching"] = strconv.FormatBool(*c.EmailSettings.EnableEmailBatching)
	props["EnablePreviewModeBanner"] = strconv.FormatBool(*c.EmailSettings.EnablePreviewModeBanner)
	props["EmailNotificationContentsType"] = *c.EmailSettings.EmailNotificationContentsType

	props["ShowEmailAddress"] = strconv.FormatBool(*c.PrivacySettings.ShowEmailAddress)
	props["ShowFullName"] = strconv.FormatBool(*c.PrivacySettings.ShowFullName)

	props["EnableFileAttachments"] = strconv.FormatBool(*c.FileSettings.EnableFileAttachments)
	props["EnablePublicLink"] = strconv.FormatBool(*c.FileSettings.EnablePublicLink)

	props["AvailableLocales"] = *c.LocalizationSettings.AvailableLocales
	props["SQLDriverName"] = *c.SqlSettings.DriverName

	props["EnableEmojiPicker"] = strconv.FormatBool(*c.ServiceSettings.EnableEmojiPicker)
	props["EnableGifPicker"] = strconv.FormatBool(*c.ServiceSettings.EnableGifPicker)
	props["GfycatApiKey"] = *c.ServiceSettings.GfycatApiKey
	props["GfycatApiSecret"] = *c.ServiceSettings.GfycatApiSecret
	props["MaxFileSize"] = strconv.FormatInt(*c.FileSettings.MaxFileSize, 10)

	props["TimeBetweenUserTypingUpdatesMilliseconds"] = strconv.FormatInt(*c.ServiceSettings.TimeBetweenUserTypingUpdatesMilliseconds, 10)
	props["EnableUserTypingMessages"] = strconv.FormatBool(*c.ServiceSettings.EnableUserTypingMessages)
	props["EnableChannelViewedMessages"] = strconv.FormatBool(*c.ServiceSettings.EnableChannelViewedMessages)

	props["RunJobs"] = strconv.FormatBool(*c.JobSettings.RunJobs)

	props["EnableEmailInvitations"] = strconv.FormatBool(*c.ServiceSettings.EnableEmailInvitations)

	props["CloudUserLimit"] = strconv.FormatInt(*c.ExperimentalSettings.CloudUserLimit, 10)

	props["EnableLegacySidebar"] = strconv.FormatBool(*c.ServiceSettings.EnableLegacySidebar)

	props["EnableReliableWebSockets"] = strconv.FormatBool(*c.ServiceSettings.EnableReliableWebSockets)

	// Set default values for all options that require a license.
	props["ExperimentalHideTownSquareinLHS"] = "false"
	props["ExperimentalTownSquareIsReadOnly"] = "false"
	props["ExperimentalEnableAuthenticationTransfer"] = "true"
	props["LdapNicknameAttributeSet"] = "false"
	props["LdapFirstNameAttributeSet"] = "false"
	props["LdapLastNameAttributeSet"] = "false"
	props["LdapPictureAttributeSet"] = "false"
	props["LdapPositionAttributeSet"] = "false"
	props["EnableCompliance"] = "false"
	props["EnableMobileFileDownload"] = "true"
	props["EnableMobileFileUpload"] = "true"
	props["SamlFirstNameAttributeSet"] = "false"
	props["SamlLastNameAttributeSet"] = "false"
	props["SamlNicknameAttributeSet"] = "false"
	props["SamlPositionAttributeSet"] = "false"
	props["EnableCluster"] = "false"
	props["EnableMetrics"] = "false"
	props["EnableBanner"] = "false"
	props["BannerText"] = ""
	props["BannerColor"] = ""
	props["BannerTextColor"] = ""
	props["AllowBannerDismissal"] = "false"
	props["EnableThemeSelection"] = "true"
	props["DefaultTheme"] = ""
	props["AllowCustomThemes"] = "true"
	props["AllowedThemes"] = ""
	props["DataRetentionEnableMessageDeletion"] = "false"
	props["DataRetentionMessageRetentionDays"] = "0"
	props["DataRetentionEnableFileDeletion"] = "false"
	props["DataRetentionFileRetentionDays"] = "0"
	props["CWSUrl"] = ""

	props["CustomUrlSchemes"] = strings.Join(c.DisplaySettings.CustomUrlSchemes, ",")
	props["IsDefaultMarketplace"] = strconv.FormatBool(*c.PluginSettings.MarketplaceUrl == model.PLUGIN_SETTINGS_DEFAULT_MARKETPLACE_URL)
	props["ExperimentalSharedChannels"] = "false"
	props["CollapsedThreads"] = *c.ServiceSettings.CollapsedThreads

	return props
}

// GenerateLimitedClientConfig renders the given configuration for an untrusted client.
func GenerateLimitedClientConfig(c *model.Config) map[string]string {
	props := make(map[string]string)

	props["Version"] = model.CurrentVersion
	props["BuildNumber"] = model.BuildNumber
	props["BuildDate"] = model.BuildDate
	props["BuildHash"] = model.BuildHash
	props["BuildHashEnterprise"] = model.BuildHashEnterprise
	props["BuildEnterpriseReady"] = model.BuildEnterpriseReady

	props["EnableBotAccountCreation"] = strconv.FormatBool(*c.ServiceSettings.EnableBotAccountCreation)
	props["EnableFile"] = strconv.FormatBool(*c.LogSettings.EnableFile)
	props["FileLevel"] = *c.LogSettings.FileLevel

	props["WebsocketURL"] = strings.TrimRight(*c.ServiceSettings.WebsocketURL, "/")
	props["WebsocketPort"] = fmt.Sprintf("%v", *c.ServiceSettings.WebsocketPort)
	props["WebsocketSecurePort"] = fmt.Sprintf("%v", *c.ServiceSettings.WebsocketSecurePort)

	props["AndroidLatestVersion"] = c.ClientRequirements.AndroidLatestVersion
	props["AndroidMinVersion"] = c.ClientRequirements.AndroidMinVersion
	props["DesktopLatestVersion"] = c.ClientRequirements.DesktopLatestVersion
	props["DesktopMinVersion"] = c.ClientRequirements.DesktopMinVersion
	props["IosLatestVersion"] = c.ClientRequirements.IosLatestVersion
	props["IosMinVersion"] = c.ClientRequirements.IosMinVersion

	props["EnableDiagnostics"] = strconv.FormatBool(*c.LogSettings.EnableDiagnostics)

	props["EnableSignUpWithEmail"] = strconv.FormatBool(*c.EmailSettings.EnableSignUpWithEmail)
	props["EnableSignInWithEmail"] = strconv.FormatBool(*c.EmailSettings.EnableSignInWithEmail)
	props["EnableSignInWithUsername"] = strconv.FormatBool(*c.EmailSettings.EnableSignInWithUsername)

	props["EmailLoginButtonColor"] = *c.EmailSettings.LoginButtonColor
	props["EmailLoginButtonBorderColor"] = *c.EmailSettings.LoginButtonBorderColor
	props["EmailLoginButtonTextColor"] = *c.EmailSettings.LoginButtonTextColor

	props["EnableSignUpWithGitLab"] = strconv.FormatBool(*c.GitLabSettings.Enable)

	props["TermsOfServiceLink"] = *c.SupportSettings.TermsOfServiceLink
	props["PrivacyPolicyLink"] = *c.SupportSettings.PrivacyPolicyLink
	props["AboutLink"] = *c.SupportSettings.AboutLink
	props["HelpLink"] = *c.SupportSettings.HelpLink
	props["ReportAProblemLink"] = *c.SupportSettings.ReportAProblemLink
	props["SupportEmail"] = *c.SupportSettings.SupportEmail
	props["EnableAskCommunityLink"] = strconv.FormatBool(*c.SupportSettings.EnableAskCommunityLink)

	props["DefaultClientLocale"] = *c.LocalizationSettings.DefaultClientLocale

	props["EnableCustomEmoji"] = strconv.FormatBool(*c.ServiceSettings.EnableCustomEmoji)
	props["AppDownloadLink"] = *c.NativeAppSettings.AppDownloadLink
	props["AndroidAppDownloadLink"] = *c.NativeAppSettings.AndroidAppDownloadLink
	props["IosAppDownloadLink"] = *c.NativeAppSettings.IosAppDownloadLink

	// props["DiagnosticId"] = telemetryID
	// props["TelemetryId"] = telemetryID
	props["DiagnosticsEnabled"] = strconv.FormatBool(*c.LogSettings.EnableDiagnostics)

	props["HasImageProxy"] = strconv.FormatBool(*c.ImageProxySettings.Enable)

	props["PluginsEnabled"] = strconv.FormatBool(*c.PluginSettings.Enable)

	props["PasswordMinimumLength"] = fmt.Sprintf("%v", *c.PasswordSettings.MinimumLength)
	props["PasswordRequireLowercase"] = strconv.FormatBool(*c.PasswordSettings.Lowercase)
	props["PasswordRequireUppercase"] = strconv.FormatBool(*c.PasswordSettings.Uppercase)
	props["PasswordRequireNumber"] = strconv.FormatBool(*c.PasswordSettings.Number)
	props["PasswordRequireSymbol"] = strconv.FormatBool(*c.PasswordSettings.Symbol)

	// Set default values for all options that require a license.
	props["EnableCustomBrand"] = "false"
	props["CustomBrandText"] = ""
	props["CustomDescriptionText"] = ""
	props["EnableLdap"] = "false"
	props["LdapLoginFieldName"] = ""
	props["LdapLoginButtonColor"] = ""
	props["LdapLoginButtonBorderColor"] = ""
	props["LdapLoginButtonTextColor"] = ""
	props["EnableSaml"] = "false"
	props["SamlLoginButtonText"] = ""
	props["SamlLoginButtonColor"] = ""
	props["SamlLoginButtonBorderColor"] = ""
	props["SamlLoginButtonTextColor"] = ""
	props["EnableSignUpWithGoogle"] = "false"
	props["EnableSignUpWithOffice365"] = "false"
	props["EnableSignUpWithOpenId"] = "false"
	props["OpenIdButtonText"] = ""
	props["OpenIdButtonColor"] = ""
	props["CWSUrl"] = ""
	props["EnableMultifactorAuthentication"] = strconv.FormatBool(*c.ServiceSettings.EnableMultifactorAuthentication)
	props["EnforceMultifactorAuthentication"] = "false"
	props["EnableGuestAccounts"] = strconv.FormatBool(*c.GuestAccountsSettings.Enable)
	props["GuestAccountsEnforceMultifactorAuthentication"] = strconv.FormatBool(*c.GuestAccountsSettings.EnforceMultifactorAuthentication)

	for key, value := range c.FeatureFlags.ToMap() {
		props["FeatureFlag"+key] = value
	}

	return props
}
