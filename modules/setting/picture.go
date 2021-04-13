package setting

import (
	"net/url"

	"github.com/sitename/sitename/modules/log"
	"strk.kbt.io/projects/go/libravatar"
)

// settings
var (
	// Picture settings
	Avatar = struct {
		Storage

		MaxWidth    int
		MaxHeight   int
		MaxFileSize int64
	}{
		MaxWidth:    4096,
		MaxHeight:   3072,
		MaxFileSize: 1048576,
	}

	GravatarSource        string
	GravatarSourceURL     *url.URL
	DisableGravatar       bool
	EnableFederatedAvatar bool
	LibravatarService     *libravatar.Libravatar
)

func newPictureService() {
	sec := Cfg.Section("picture")

	avatarSec := Cfg.Section("avatar")
	storageType := sec.Key("AVATAR_STORAGE_TYPE").MustString("")
	// Specifically default PATH to AVATAR_UPLOAD_PATH
	avatarSec.Key("PATH").MustString(
		sec.Key("AVATAR_UPLOAD_PATH").String())

	Avatar.Storage = getStorage("avatars", storageType, avatarSec)

	Avatar.MaxWidth = sec.Key("AVATAR_MAX_WIDTH").MustInt(4096)
	Avatar.MaxHeight = sec.Key("AVATAR_MAX_HEIGHT").MustInt(3072)
	Avatar.MaxFileSize = sec.Key("AVATAR_MAX_FILE_SIZE").MustInt64(1048576)

	switch source := sec.Key("GRAVATAR_SOURCE").MustString("gravatar"); source {
	case "duoshuo":
		GravatarSource = "http://gravatar.duoshuo.com/avatar/"
	case "gravatar":
		GravatarSource = "https://secure.gravatar.com/avatar/"
	case "libravatar":
		GravatarSource = "https://seccdn.libravatar.org/avatar/"
	default:
		GravatarSource = source
	}
	DisableGravatar = sec.Key("DISABLE_GRAVATAR").MustBool()
	EnableFederatedAvatar = sec.Key("ENABLE_FEDERATED_AVATAR").MustBool(!InstallLock)
	if OfflineMode {
		DisableGravatar = true
		EnableFederatedAvatar = false
	}
	if DisableGravatar {
		EnableFederatedAvatar = false
	}
	if EnableFederatedAvatar || !DisableGravatar {
		var err error
		GravatarSourceURL, err = url.Parse(GravatarSource)
		if err != nil {
			log.Fatal("Failed to parse Gravatar URL(%s): %v",
				GravatarSource, err)
		}
	}

	if EnableFederatedAvatar {
		LibravatarService = libravatar.New()
		if GravatarSourceURL.Scheme == "https" {
			LibravatarService.SetUseHTTPS(true)
			LibravatarService.SetSecureFallbackHost(GravatarSourceURL.Host)
		} else {
			LibravatarService.SetUseHTTPS(false)
			LibravatarService.SetFallbackHost(GravatarSourceURL.Host)
		}
	}
}
