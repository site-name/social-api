package context

import (
	"sync"

	"github.com/sitename/sitename/modules/cache"
	"github.com/sitename/sitename/modules/setting"

	"gitea.com/go-chi/captcha"
)

var imageCaptchaOnce sync.Once
var cpt *captcha.Captcha

// GetImageCaptcha returns global image captcha
func GetImageCaptcha() *captcha.Captcha {
	imageCaptchaOnce.Do(func() {
		cpt = captcha.NewCaptcha(captcha.Options{
			SubURL: setting.AppSubURL,
		})
		cpt.Store = cache.GetCache()
	})
	return cpt
}
