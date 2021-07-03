package file

import (
	"fmt"
	"io"
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/imaging"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/slog"
)

type AppFile struct {
	app.AppIface
	// These are used to prevent concurrent upload requests
	// for a given upload session which could cause inconsistencies
	// and data corruption.
	uploadLockMapMut sync.Mutex
	uploadLockMap    map[string]bool
}

func init() {
	app.RegisterFileApp(func(a app.AppIface) sub_app_iface.FileApp {
		fa := &AppFile{
			AppIface:      a,
			uploadLockMap: map[string]bool{},
		}

		// test file backend connection:
		backend, appErr := fa.FileBackend()
		if appErr != nil {
			slog.Error("Problem with file storage settings", slog.Err(appErr))
		} else {
			nErr := backend.TestConnection()
			if nErr != nil {
				if _, ok := nErr.(*filestore.S3FileBackendNoBucketError); ok {
					nErr = backend.(*filestore.S3FileBackend).MakeBucket()
				}
				if nErr != nil {
					slog.Error("Problem with file storage settings", slog.Err(nErr))
				}
			}
		}

		return fa
	})
}

func checkImageResolutionLimit(w, h int) error {
	// This casting is done to prevent overflow on 32 bit systems (not needed
	// in 64 bits systems because images can't have more than 32 bits height or
	// width)
	imageRes := int64(w) * int64(h)
	if imageRes > maxImageRes {
		return fmt.Errorf("image resolution is too high: %d, max allowed is %d", imageRes, maxImageRes)
	}

	return nil
}

func checkImageLimits(imageData io.Reader) error {
	w, h, err := imaging.GetDimensions(imageData)
	if err != nil {
		return fmt.Errorf("failed to get image dimensions: %w", err)
	}

	return checkImageResolutionLimit(w, h)
}
