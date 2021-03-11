package models

import (
	"net/url"
	"path"
	"strings"

	"github.com/sitename/sitename/modules/base"
	"github.com/sitename/sitename/modules/cache"
	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/modules/setting"
)

// EmailHash represents a pre-generated hash map
type EmailHash struct {
	Hash  string `xorm:"pk VARCHAR(32)"`
	Email string `xorm:"UNIQUE NOT NULL"`
}

// DefaultAvatarLink the default avatar link
func DefaultAvatarLink() string {
	u, err := url.Parse(setting.AppSubURL)
	if err != nil {
		log.Error("GetUserByEmail: %v", err)
		return ""
	}

	u.Path = path.Join(u.Path, "/img/avatar_default.png")
	return u.String()
}

// DefaultAvatarSize is a sentinel value for the default avatar size, as
// determined by the avatar-hosting service.
const DefaultAvatarSize = -1

// DefaultAvatarPixelSize is the default size in pixels of a rendered avatar
const DefaultAvatarPixelSize = 28

// AvatarRenderedSizeFactor is the factor by which the default size is increased for finer rendering
const AvatarRenderedSizeFactor = 2

// HashEmail hashes email address to MD5 string.
// https://en.gravatar.com/site/implement/hash/
func HashEmail(email string) string {
	return base.EncodeMD5(strings.ToLower(strings.TrimSpace(email)))
}

// GetEmailForHash converts a provided md5sum to the email
func GetEmailForHash(md5sum string) (string, error) {
	return cache.GetString("Avatar:"+md5sum, func() (string, error) {
		emailHash := EmailHash{
			Hash: strings.ToLower(strings.TrimSpace(md5sum)),
		}

		_, err := x.Get(&emailHash)
		return emailHash.Email, err
	})
}
