package file

import (
	"sync"

	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/sub_app_iface"
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

type ServiceFileConfig struct {
	Server *app.Server
}

func NewServiceFile(config *ServiceFileConfig) (sub_app_iface.FileService, error) {
	fa := &ServiceFile{
		srv:           config.Server,
		uploadLockMap: map[string]bool{},
	}

	// test file backend connection:
	backend, appErr := fa.FileBackend()
	if appErr != nil {
		return nil, appErr
	} else {
		nErr := backend.TestConnection()
		if nErr != nil {
			if _, ok := nErr.(*filestore.S3FileBackendNoBucketError); ok {
				nErr = backend.(*filestore.S3FileBackend).MakeBucket()
			}
			if nErr != nil {
				return nil, nErr
			}
		}
	}

	return fa, nil
}
