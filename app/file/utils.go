package file

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/sitename/sitename/app/imaging"
	"github.com/sitename/sitename/model/file"
)

func checkImageResolutionLimit(w, h int, maxRes int64) error {
	// This casting is done to prevent overflow on 32 bit systems (not needed
	// in 64 bits systems because images can't have more than 32 bits height or
	// width)
	imageRes := int64(w) * int64(h)
	if imageRes > maxRes {
		return fmt.Errorf("image resolution is too high: %d, max allowed is %d", imageRes, maxRes)
	}

	return nil
}

// CheckImageLimits
func CheckImageLimits(imageData io.Reader, maxRes int64) error {
	w, h, err := imaging.GetDimensions(imageData)
	if err != nil {
		return fmt.Errorf("failed to get image dimensions: %w", err)
	}

	return checkImageResolutionLimit(w, h, maxRes)
}

func (a *ServiceFile) GeneratePublicLink(siteURL string, info *file.FileInfo) string {
	hash := GeneratePublicLinkHash(info.Id, *a.srv.Config().FileSettings.PublicLinkSalt)
	return fmt.Sprintf("%s/files/%v/public?h=%s", siteURL, info.Id, hash)
}

func GeneratePublicLinkHash(fileID, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(salt))
	hash.Write([]byte(fileID))

	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}
