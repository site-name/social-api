package models

import (
	"fmt"
	"image/png"
	"io"

	"github.com/sitename/sitename/modules/avatar"
	"github.com/sitename/sitename/modules/log"
	"github.com/sitename/sitename/modules/storage"
)

// CustomAvatarRelativePath returns user custom avatar relative path.
func (u *User) CustomAvatarRelativePath() string {
	return u.Avatar
}

// GenerateRandomAvatar generates a random avatar for user.
func (u *User) GenerateRandomAvatar() error {
	return u.generateRandomAvatar(x)
}

func (u *User) generateRandomAvatar(e Engine) error {
	seed := u.Email
	if len(seed) == 0 {
		seed = u.Name
	}

	img, err := avatar.RandomImage([]byte(seed))
	if err != nil {
		return fmt.Errorf("RandomImage: %v", err)
	}

	u.Avatar = HashEmail(seed)

	if err := storage.SaveFrom(storage.Avatars, u.CustomAvatarRelativePath(), func(w io.Writer) error {
		if err := png.Encode(w, img); err != nil {
			log.Error("encode: %v", err)
		}
		return err
	}); err != nil {
		return fmt.Errorf("Failed to create dir %s: %v", u.CustomAvatarRelativePath(), err)
	}

	if _, err := e.ID(u.ID).Cols("avatar").Update(u); err != nil {
		return err
	}

	log.Info("New random avatar created: %d", u.ID)
	return nil
}
