/*
NOTE: This package is initialized during server startup (modules/imports does that)
so the init() function get the chance to register a function to create `ServiceAccount`
*/
package file

import (
	"io"
	"net/http"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/imaging"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/services/docextractor"
)

const (
	imageThumbnailWidth        = 120
	imageThumbnailHeight       = 100
	imagePreviewWidth          = 1920
	miniPreviewImageWidth      = 16
	miniPreviewImageHeight     = 16
	jpegEncQuality             = 90
	maxUploadInitialBufferSize = 1024 * 1024 // 1MB
	maxContentExtractionSize   = 1024 * 1024 // 1MB
)

type ServiceFile struct {
	srv *app.Server
	// These are used to prevent concurrent upload requests
	// for a given upload session which could cause inconsistencies
	// and data corruption.
	uploadLockMapMut sync.Mutex
	uploadLockMap    map[string]bool

	imgDecoder *imaging.Decoder
	imgEncoder *imaging.Encoder
}

func init() {
	app.RegisterService(func(s *app.Server) error {
		service := &ServiceFile{
			srv:           s,
			uploadLockMap: map[string]bool{},
		}

		var err error
		service.imgDecoder, err = imaging.NewDecoder(imaging.DecoderOptions{
			ConcurrencyLevel: runtime.NumCPU(),
		})
		if err != nil {
			return err
		}

		service.imgEncoder, err = imaging.NewEncoder(imaging.EncoderOptions{
			ConcurrencyLevel: runtime.NumCPU(),
		})
		if err != nil {
			return err
		}

		// test file backend connection:
		backend, appErr := service.FileBackend()
		if appErr != nil {
			return appErr
		}

		nErr := backend.TestConnection()
		if nErr != nil {
			if _, ok := nErr.(*filestore.S3FileBackendNoBucketError); ok {
				nErr = backend.(*filestore.S3FileBackend).MakeBucket()
			}
			if nErr != nil {
				return nErr
			}
		}

		s.File = service
		return nil
	})
}

// ImageEncoder returns image encoder
func (s *ServiceFile) ImageEncoder() *imaging.Encoder {
	return s.imgEncoder
}

// ImageDecoder retutns image encoder
func (s *ServiceFile) ImageDecoder() *imaging.Decoder {
	return s.imgDecoder
}

var oldFilenameMatchExp *regexp.Regexp = regexp.MustCompile(`^\/([a-z\d]{26})\/([a-z\d]{26})\/([a-z\d]{26})\/([^\/]+)$`)

// // Parse the path from the Filename of the form /{channelID}/{userID}/{uid}/{nameWithExtension}
func parseOldFilenames(filenames []string, channelID, userID string) [][]string {
	parsed := [][]string{}
	for _, filename := range filenames {
		matches := oldFilenameMatchExp.FindStringSubmatch(filename)
		if len(matches) != 5 {
			slog.Error("Failed to parse old Filename", slog.String("filename", filename))
			continue
		}
		if matches[1] != channelID {
			slog.Error("ChannelId in Filename does not match", slog.String("channel_id", channelID), slog.String("matched", matches[1]))
		} else if matches[2] != userID {
			slog.Error("UserId in Filename does not match", slog.String("user_id", userID), slog.String("matched", matches[2]))
		} else {
			parsed = append(parsed, matches[1:])
		}
	}
	return parsed
}

// FileBackend returns filebackend of the system
func (a *ServiceFile) FileBackend() (filestore.FileBackend, *model_helper.AppError) {
	backend, err := filestore.NewFileBackend(a.srv.Config().FileSettings.ToFileBackendSettings(true))
	if err != nil {
		return nil, model_helper.NewAppError("FileBackend", "app.file.no_driver.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return backend, nil
}

func (a *ServiceFile) CheckMandatoryS3Fields(settings *model.FileSettings) *model_helper.AppError {
	fileBackendSettings := settings.ToFileBackendSettings(false)
	err := fileBackendSettings.CheckMandatoryS3Fields()
	if err != nil {
		return model_helper.NewAppError("CheckMandatoryS3Fields", "api.admin.test_s3.missing_s3_bucket", nil, err.Error(), http.StatusBadRequest)
	}
	return nil
}

// convert filebackend connection error to system's standard app error
func connectionTestErrorToAppError(connTestErr error) *model_helper.AppError {
	switch err := connTestErr.(type) {
	case *filestore.S3FileBackendAuthError:
		return model_helper.NewAppError("TestConnection", "api.file.test_connection_s3_auth.app_error", nil, err.Error(), http.StatusInternalServerError)
	case *filestore.S3FileBackendNoBucketError:
		return model_helper.NewAppError("TestConnection", "api.file.test_connection_s3_bucket_does_not_exist.app_error", nil, err.Error(), http.StatusInternalServerError)
	default:
		return model_helper.NewAppError("TestConnection", "api.file.test_connection.app_error", nil, connTestErr.Error(), http.StatusInternalServerError)
	}
}

// TestFileStoreConnection test if connection to file backend server is good
func (a *ServiceFile) TestFileStoreConnection() *model_helper.AppError {
	backend, err := a.FileBackend()
	if err != nil {
		return err
	}
	nErr := backend.TestConnection()
	if nErr != nil {
		return connectionTestErrorToAppError(nErr)
	}
	return nil
}

// TestFileStoreConnectionWithConfig test file backend connection with config
func (a *ServiceFile) TestFileStoreConnectionWithConfig(settings *model.FileSettings) *model_helper.AppError {
	backend, err := filestore.NewFileBackend(settings.ToFileBackendSettings(true))
	if err != nil {
		return model_helper.NewAppError("FileBackend", "api.file.no_driver.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	nErr := backend.TestConnection()
	if nErr != nil {
		return connectionTestErrorToAppError(nErr)
	}
	return nil
}

// ReadFile read file content from given path
func (a *ServiceFile) ReadFile(path string) ([]byte, *model_helper.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return nil, err
	}
	result, nErr := backend.ReadFile(path)
	if nErr != nil {
		return nil, model_helper.NewAppError("ReadFile", "api.file.read_file.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

// Caller must close the first return value
func (a *ServiceFile) FileReader(path string) (filestore.ReadCloseSeeker, *model_helper.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return nil, err
	}
	result, nErr := backend.Reader(path)
	if nErr != nil {
		return nil, model_helper.NewAppError("FileReader", "api.file.file_reader.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

// FileExists checks if given path exists
func (a *ServiceFile) FileExists(path string) (bool, *model_helper.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return false, err
	}
	result, nErr := backend.FileExists(path)
	if nErr != nil {
		return false, model_helper.NewAppError("FileExists", "api.file.file_exists.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

// FileSize checks size of given path
func (a *ServiceFile) FileSize(path string) (int64, *model_helper.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return 0, err
	}
	size, nErr := backend.FileSize(path)
	if nErr != nil {
		return 0, model_helper.NewAppError("FileSize", "api.file.file_size.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return size, nil
}

// FileModTime get last modification time of given path
func (a *ServiceFile) FileModTime(path string) (time.Time, *model_helper.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return time.Time{}, err
	}
	modTime, nErr := backend.FileModTime(path)
	if nErr != nil {
		return time.Time{}, model_helper.NewAppError("FileModTime", "api.file.file_mod_time.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}

	return modTime, nil
}

// MoveFile moves file from given oldPath to newPath
func (a *ServiceFile) MoveFile(oldPath, newPath string) *model_helper.AppError {
	backend, err := a.FileBackend()
	if err != nil {
		return err
	}
	nErr := backend.MoveFile(oldPath, newPath)
	if nErr != nil {
		return model_helper.NewAppError("MoveFile", "api.file.move_file.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (a *ServiceFile) WriteFile(fr io.Reader, path string) (int64, *model_helper.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return 0, err
	}

	result, nErr := backend.WriteFile(fr, path)
	if nErr != nil {
		return result, model_helper.NewAppError("WriteFile", "api.file.write_file.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

func (a *ServiceFile) AppendFile(fr io.Reader, path string) (int64, *model_helper.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return 0, err
	}

	result, nErr := backend.AppendFile(fr, path)
	if nErr != nil {
		return result, model_helper.NewAppError("AppendFile", "api.file.append_file.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return result, nil
}

func (a *ServiceFile) RemoveFile(path string) *model_helper.AppError {
	backend, err := a.FileBackend()
	if err != nil {
		return err
	}
	nErr := backend.RemoveFile(path)
	if nErr != nil {
		return model_helper.NewAppError("RemoveFile", "api.file.remove_file.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (a *ServiceFile) ListDirectory(path string) ([]string, *model_helper.AppError) {
	backend, err := a.FileBackend()
	if err != nil {
		return nil, err
	}
	paths, nErr := backend.ListDirectory(path)
	if nErr != nil {
		return nil, model_helper.NewAppError("ListDirectory", "api.file.list_directory.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}

	return paths, nil
}

func (a *ServiceFile) RemoveDirectory(path string) *model_helper.AppError {
	backend, err := a.FileBackend()
	if err != nil {
		return err
	}
	nErr := backend.RemoveDirectory(path)
	if nErr != nil {
		return model_helper.NewAppError("RemoveDirectory", "api.file.remove_directory.app_error", nil, nErr.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (a *ServiceFile) ExtractContentFromFileInfo(fileInfo *model.FileInfo) error {
	file, aerr := a.FileReader(fileInfo.Path)
	if aerr != nil {
		return errors.Wrap(aerr, "failed to open file for extract file content")
	}
	defer file.Close()

	text, err := docextractor.Extract(fileInfo.Name, file, docextractor.ExtractSettings{
		ArchiveRecursion: *a.srv.Config().FileSettings.ArchiveRecursion,
	})
	if err != nil {
		return errors.Wrap(err, "failed to extract file content")
	}
	if text != "" {
		if len(text) > maxContentExtractionSize {
			text = text[0:maxContentExtractionSize]
		}
		if _, storeErr := a.srv.Store.FileInfo().Upsert(&model.FileInfo{Id: fileInfo.Id, Content: text}); storeErr != nil {
			return errors.Wrap(storeErr, "failed to save the extracted file content")
		}
		_, storeErr := a.srv.Store.FileInfo().Get(fileInfo.Id, false)
		if storeErr != nil {
			slog.Warn("failed to invalidate the fileInfo cache.", slog.Err(storeErr), slog.String("file_info_id", fileInfo.Id))
		} else {
			// a.srv.Store.FileInfo().InvalidateFileInfosForPostCache()
			slog.Warn("This flow is not implemented", slog.String("function", "generateMiniPreview"))
		}
	}
	return nil
}
