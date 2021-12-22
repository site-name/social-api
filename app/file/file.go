/*
	NOTE: This package is initialized during server startup (modules/imports does that)
	so the init() function get the chance to register a function to create `ServiceAccount`
*/
package file

import (
	"net/http"
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/filestore"
)

type ServiceFile struct {
	srv *app.Server
	// These are used to prevent concurrent upload requests
	// for a given upload session which could cause inconsistencies
	// and data corruption.
	uploadLockMapMut sync.Mutex
	uploadLockMap    map[string]bool
}

func init() {
	app.RegisterFileService(func(s *app.Server) (sub_app_iface.FileService, error) {
		service := &ServiceFile{
			srv:           s,
			uploadLockMap: map[string]bool{},
		}

		// test file backend connection:
		backend, appErr := service.FileBackend()
		if appErr != nil {
			return nil, appErr
		}

		nErr := backend.TestConnection()
		if nErr != nil {
			if _, ok := nErr.(*filestore.S3FileBackendNoBucketError); ok {
				nErr = backend.(*filestore.S3FileBackend).MakeBucket()
			}
			if nErr != nil {
				return nil, nErr
			}
		}

		return service, nil
	})
}

// FileBackend returns filebackend of the system
func (a *ServiceFile) FileBackend() (filestore.FileBackend, *model.AppError) {
	backend, err := filestore.NewFileBackend(a.srv.Config().FileSettings.ToFileBackendSettings(true))
	if err != nil {
		return nil, model.NewAppError("FileBackend", "app.file.no_driver.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return backend, nil
}
