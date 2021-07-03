package file

import (
	"sync"

	"github.com/sitename/sitename/app"
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
