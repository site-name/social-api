package file

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"image/gif"
	"io"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/sitename/sitename/app"
	"github.com/sitename/sitename/app/imaging"
	"github.com/sitename/sitename/app/request"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model_helper"
	"github.com/sitename/sitename/modules/filestore"
	"github.com/sitename/sitename/modules/model_types"
	"github.com/sitename/sitename/modules/plugin"
	"github.com/sitename/sitename/modules/slog"
	"github.com/sitename/sitename/store"
	"github.com/volatiletech/null/v8"
)

// func (a *ServiceFile) getInfoForFilename(post *model.Post, teamID, channelID, userID, oldId, filename string) *model.FileInfo {
// 	name, _ := url.QueryUnescape(filename)
// 	pathPrefix := fmt.Sprintf("teams/%s/channels/%s/users/%s/%s/", teamID, channelID, userID, oldId)
// 	path := pathPrefix + name

// 	// Open the file and populate the fields of the FileInfo
// 	data, err := a.ReadFile(path)
// 	if err != nil {
// 		slog.Error(
// 			"File not found when migrating post to use FileInfos",
// 			slog.String("post_id", post.Id),
// 			slog.String("filename", filename),
// 			slog.String("path", path),
// 			slog.Err(err),
// 		)
// 		return nil
// 	}

// 	info, err := model.GetInfoForBytes(name, bytes.NewReader(data), len(data))
// 	if err != nil {
// 		slog.Warn(
// 			"Unable to fully decode file info when migrating post to use FileInfos",
// 			slog.String("post_id", post.Id),
// 			slog.String("filename", filename),
// 			slog.Err(err),
// 		)
// 	}

// 	// Generate a new ID because with the old system, you could very rarely get multiple posts referencing the same file
// 	info.Id = model.NewId()
// 	info.CreatorId = post.UserId
// 	info.PostId = post.Id
// 	info.CreateAt = post.CreateAt
// 	info.UpdateAt = post.UpdateAt
// 	info.Path = path

// 	if info.IsImage() {
// 		nameWithoutExtension := name[:strings.LastIndex(name, ".")]
// 		info.PreviewPath = pathPrefix + nameWithoutExtension + "_preview.jpg"
// 		info.ThumbnailPath = pathPrefix + nameWithoutExtension + "_thumb.jpg"
// 	}

// 	return info
// }

// func (a *ServiceFile) findTeamIdForFilename(post *model.Post, id, filename string) string {
// 	name, _ := url.QueryUnescape(filename)

// 	// This post is in a direct channel so we need to figure out what team the files are stored under.
// 	teams, err := a.srv.Store.Team().GetTeamsByUserId(post.UserId)
// 	if err != nil {
// 		slog.Error("Unable to get teams when migrating post to use FileInfo", slog.Err(err), slog.String("post_id", post.Id))
// 		return ""
// 	}

// 	if len(teams) == 1 {
// 		// The user has only one team so the post must've been sent from it
// 		return teams[0].Id
// 	}

// 	for _, team := range teams {
// 		path := fmt.Sprintf("teams/%s/channels/%s/users/%s/%s/%s", team.Id, post.ChannelId, post.UserId, id, name)
// 		if ok, err := a.FileExists(path); ok && err == nil {
// 			// Found the team that this file was posted from
// 			return team.Id
// 		}
// 	}

// 	return ""
// }

// var fileMigrationLock sync.Mutex

// // Creates and stores FileInfos for a post created before the FileInfos table existed.
// func (a *ServiceFile) MigrateFilenamesToFileInfos(post *model.Post) []*model.FileInfo {
// 	if len(post.Filenames) == 0 {
// 		slog.Warn("Unable to migrate post to use FileInfos with an empty Filenames field", slog.String("post_id", post.Id))
// 		return []*model.FileInfo{}
// 	}

// 	channel, errCh := a.srv.Store.Channel().Get(post.ChannelId, true)
// 	// There's a weird bug that rarely happens where a post ends up with duplicate Filenames so remove those
// 	filenames := util.Dedup(post.Filenames)
// 	if errCh != nil {
// 		slog.Error(
// 			"Unable to get channel when migrating post to use FileInfos",
// 			slog.String("post_id", post.Id),
// 			slog.String("channel_id", post.ChannelId),
// 			slog.Err(errCh),
// 		)
// 		return []*model.FileInfo{}
// 	}

// 	// Parse and validate filenames before further processing
// 	parsedFilenames := parseOldFilenames(filenames, post.ChannelId, post.UserId)

// 	if len(parsedFilenames) == 0 {
// 		slog.Error("Unable to parse filenames")
// 		return []*model.FileInfo{}
// 	}

// 	// Find the team that was used to make this post since its part of the file path that isn't saved in the Filename
// 	var teamID string
// 	if channel.TeamId == "" {
// 		// This post was made in a cross-team DM channel, so we need to find where its files were saved
// 		teamID = a.findTeamIdForFilename(post, parsedFilenames[0][2], parsedFilenames[0][3])
// 	} else {
// 		teamID = channel.TeamId
// 	}

// 	// Create FileInfo objects for this post
// 	infos := make([]*model.FileInfo, 0, len(filenames))
// 	if teamID == "" {
// 		slog.Error(
// 			"Unable to find team id for files when migrating post to use FileInfos",
// 			slog.String("filenames", strings.Join(filenames, ",")),
// 			slog.String("post_id", post.Id),
// 		)
// 	} else {
// 		for _, parsed := range parsedFilenames {
// 			info := a.getInfoForFilename(post, teamID, parsed[0], parsed[1], parsed[2], parsed[3])
// 			if info == nil {
// 				continue
// 			}

// 			infos = append(infos, info)
// 		}
// 	}

// 	// Lock to prevent only one migration thread from trying to update the post at once, preventing duplicate FileInfos from being created
// 	fileMigrationLock.Lock()
// 	defer fileMigrationLock.Unlock()

// 	result, nErr := a.srv.Store.Post().Get(context.Background(), post.Id, false, false, false, "")
// 	if nErr != nil {
// 		slog.Error("Unable to get post when migrating post to use FileInfos", slog.Err(nErr), slog.String("post_id", post.Id))
// 		return []*model.FileInfo{}
// 	}

// 	if newPost := result.Posts[post.Id]; len(newPost.Filenames) != len(post.Filenames) {
// 		// Another thread has already created FileInfos for this post, so just return those
// 		var fileInfos []*model.FileInfo
// 		fileInfos, nErr = a.srv.Store.FileInfo().GetForPost(post.Id, true, false, false)
// 		if nErr != nil {
// 			slog.Error("Unable to get FileInfos for migrated post", slog.Err(nErr), slog.String("post_id", post.Id))
// 			return []*model.FileInfo{}
// 		}

// 		slog.Debug("Post already migrated to use FileInfos", slog.String("post_id", post.Id))
// 		return fileInfos
// 	}

// 	slog.Debug("Migrating post to use FileInfos", slog.String("post_id", post.Id))

// 	savedInfos := make([]*model.FileInfo, 0, len(infos))
// 	fileIDs := make([]string, 0, len(filenames))
// 	for _, info := range infos {
// 		if _, nErr = a.srv.Store.FileInfo().Save(info); nErr != nil {
// 			slog.Error(
// 				"Unable to save file info when migrating post to use FileInfos",
// 				slog.String("post_id", post.Id),
// 				slog.String("file_info_id", info.Id),
// 				slog.String("file_info_path", info.Path),
// 				slog.Err(nErr),
// 			)
// 			continue
// 		}

// 		savedInfos = append(savedInfos, info)
// 		fileIDs = append(fileIDs, info.Id)
// 	}

// 	// Copy and save the updated post
// 	newPost := post.Clone()

// 	newPost.Filenames = []string{}
// 	newPost.FileIds = fileIDs

// 	// Update Posts to clear Filenames and set FileIds
// 	if _, nErr = a.srv.Store.Post().Update(newPost, post); nErr != nil {
// 		slog.Error(
// 			"Unable to save migrated post when migrating to use FileInfos",
// 			slog.String("new_file_ids", strings.Join(newPost.FileIds, ",")),
// 			slog.String("old_filenames", strings.Join(post.Filenames, ",")),
// 			slog.String("post_id", post.Id),
// 			slog.Err(nErr),
// 		)
// 		return []*model.FileInfo{}
// 	}
// 	return savedInfos
// }

// func (a *ServiceFile) UploadMultipartFiles(c *request.Context, teamID string, channelID string, userID string, fileHeaders []*multipart.FileHeader, clientIds []string, now time.Time) (*file.FileUploadResponse, *model_helper.AppError) {
// 	files := make([]io.ReadCloser, len(fileHeaders))
// 	filenames := make([]string, len(fileHeaders))

// 	for i, fileHeader := range fileHeaders {
// 		file, fileErr := fileHeader.Open()
// 		if fileErr != nil {
// 			return nil, model_helper.NewAppError(
// 				"UploadFiles",
// 				"api.file.upload_file.read_request.app_error",
// 				map[string]any{
// 					"Filename": fileHeader.Filename,
// 				},
// 				fileErr.Error(),
// 				http.StatusBadRequest,
// 			)
// 		}

// 		// Will be closed after UploadFiles returns
// 		defer file.Close()

// 		files[i] = file
// 		filenames[i] = fileHeader.Filename
// 	}

// 	return a.UploadFiles(c, teamID, channelID, userID, files, filenames, clientIds, now)
// }

// Uploads some files to the given team and channel as the given user. files and filenames should have
// the same length. clientIds should either not be provided or have the same length as files and filenames.
// The provided files should be closed by the caller so that they are not leaked.
// func (a *ServiceFile) UploadFiles(c *request.Context, teamID string, channelID string, userID string, files []io.ReadCloser, filenames []string, clientIds []string, now time.Time) (*file.FileUploadResponse, *model_helper.AppError) {
// 	if *a.srv.Config().FileSettings.DriverName == "" {
// 		return nil, model_helper.NewAppError("UploadFiles", "api.file.upload_file.storage.app_error", nil, "", http.StatusNotImplemented)
// 	}

// 	if len(filenames) != len(files) || (len(clientIds) > 0 && len(clientIds) != len(files)) {
// 		return nil, model_helper.NewAppError("UploadFiles", "api.file.upload_file.incorrect_number_of_files.app_error", nil, "", http.StatusBadRequest)
// 	}

// 	resStruct := &file.FileUploadResponse{
// 		FileInfos: []*file.FileInfo{},
// 		ClientIds: []string{},
// 	}

// 	previewPathList := []string{}
// 	thumbnailPathList := []string{}
// 	imageDataList := [][]byte{}

// 	for i, file := range files {
// 		buf := bytes.NewBuffer(nil)
// 		io.Copy(buf, file)
// 		data := buf.Bytes()

// 		info, data, err := a.DoUploadFileExpectModification(c, now, teamID, channelID, userID, filenames[i], data)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if info.PreviewPath != "" || info.ThumbnailPath != "" {
// 			previewPathList = append(previewPathList, info.PreviewPath)
// 			thumbnailPathList = append(thumbnailPathList, info.ThumbnailPath)
// 			imageDataList = append(imageDataList, data)
// 		}

// 		resStruct.FileInfos = append(resStruct.FileInfos, info)

// 		if len(clientIds) > 0 {
// 			resStruct.ClientIds = append(resStruct.ClientIds, clientIds[i])
// 		}
// 	}

// 	a.HandleImages(previewPathList, thumbnailPathList, imageDataList)

// 	return resStruct, nil
// }

// UploadFile uploads a single file in form of a completely constructed byte array for a channel.
// func (a *ServiceFile) UploadFile(c *request.Context, data []byte, channelID string, filename string) (*file.FileInfo, *model_helper.AppError) {
// 	_, err := a.GetChannel(channelID)
// 	if err != nil && channelID != "" {
// 		return nil, model_helper.NewAppError("UploadFile", "api.file.upload_file.incorrect_channelId.app_error", map[string]any{"channelId": channelID}, "", http.StatusBadRequest)
// 	}

// 	info, _, appError := a.DoUploadFileExpectModification(c, time.Now(), "noteam", channelID, "nouser", filename, data)
// 	if appError != nil {
// 		return nil, appError
// 	}

// 	if info.PreviewPath != "" || info.ThumbnailPath != "" {
// 		previewPathList := []string{info.PreviewPath}
// 		thumbnailPathList := []string{info.ThumbnailPath}
// 		imageDataList := [][]byte{data}

// 		a.HandleImages(previewPathList, thumbnailPathList, imageDataList)
// 	}

// 	return info, nil
// }

func (a *ServiceFile) DoUploadFile(c *request.Context, now time.Time, rawTeamId string, rawChannelId string, rawUserId string, rawFilename string, data []byte) (*model.FileInfo, *model_helper.AppError) {
	info, _, err := a.DoUploadFileExpectModification(c, now, rawTeamId, rawChannelId, rawUserId, rawFilename, data)
	return info, err
}

// func UploadFileSetTeamId(teamID string) func(t *UploadFileTask) {
// 	return func(t *UploadFileTask) {
// 		t.TeamId = filepath.Base(teamID)
// 	}
// }

// func UploadFileSetUserId(userID string) func(t *UploadFileTask) {
// 	return func(t *UploadFileTask) {
// 		t.UserId = filepath.Base(userID)
// 	}
// }

// func UploadFileSetTimestamp(timestamp time.Time) func(t *UploadFileTask) {
// 	return func(t *UploadFileTask) {
// 		t.Timestamp = timestamp
// 	}
// }

// func UploadFileSetContentLength(contentLength int64) func(t *UploadFileTask) {
// 	return func(t *UploadFileTask) {
// 		t.ContentLength = contentLength
// 	}
// }

// func UploadFileSetClientId(clientId string) func(t *UploadFileTask) {
// 	return func(t *UploadFileTask) {
// 		t.ClientId = clientId
// 	}
// }

// func UploadFileSetRaw() func(t *UploadFileTask) {
// 	return func(t *UploadFileTask) {
// 		t.Raw = true
// 	}
// }

type UploadFileTask struct {
	Name          string // File name
	UserId        string
	Timestamp     time.Time // Time stamp to use when creating the file.
	ContentLength int64     // The value of the Content-Length http header, when available.
	Input         io.Reader // The file data stream.
	ClientId      string    // An optional, client-assigned Id field.
	Raw           bool      // If Raw, do not execute special processing for images, just upload the file.  Plugins are still invoked.

	buf          *bytes.Buffer
	limit        int64
	limitedInput io.Reader
	teeInput     io.Reader
	fileinfo     *model.FileInfo
	maxFileSize  int64
	maxImageRes  int64

	decoded          image.Image // Cached image data that (may) get initialized in preprocessImage and is used in postprocessImage
	imageType        string
	imageOrientation int

	writeFile      func(io.Reader, string) (int64, *model_helper.AppError)
	saveToDatabase func(model.FileInfo) (*model.FileInfo, error)

	imgDecoder         *imaging.Decoder
	imgEncoder         *imaging.Encoder
	pluginsEnvironment *plugin.Environment
}

func (t *UploadFileTask) init(a *ServiceFile) {
	t.buf = &bytes.Buffer{}
	if t.ContentLength > 0 {
		t.limit = t.ContentLength
	} else {
		t.limit = t.maxFileSize
	}

	if t.ContentLength > 0 && t.ContentLength < maxUploadInitialBufferSize {
		t.buf.Grow(int(t.ContentLength))
	} else {
		t.buf.Grow(maxUploadInitialBufferSize)
	}

	model_helper.GetMillis()

	t.fileinfo = model_helper.NewFileInfo(filepath.Base(t.Name))
	t.fileinfo.ID = model_helper.NewId()
	t.fileinfo.CreatorID = t.UserId
	t.fileinfo.CreatedAt = t.Timestamp.UnixNano() / int64(time.Millisecond)
	t.fileinfo.Path = t.pathPrefix() + t.Name

	t.limitedInput = &io.LimitedReader{
		R: t.Input,
		N: t.limit + 1,
	}
	t.teeInput = io.TeeReader(t.limitedInput, t.buf)

	t.pluginsEnvironment, _ = a.srv.Plugin.GetPluginsEnvironment()
	t.writeFile = a.WriteFile
	t.saveToDatabase = a.srv.Store.FileInfo().Upsert
}

// UploadFileX uploads a single file as specified in t. It applies the upload
// constraints, executes plugins and image processing logic as needed. It
// returns a filled-out FileInfo and an optional error. A plugin may reject the
// upload, returning a rejection error. In this case FileInfo would have
// contained the last "good" FileInfo before the execution of that plugin.
func (a *ServiceFile) UploadFileX(c *request.Context, channelID, name string, input io.Reader /* after are optional arguments */, userID *string, timestamp *time.Time, contentLength *int64, clientID *string, raw *bool) (*model.FileInfo, *model_helper.AppError) {
	t := &UploadFileTask{
		// ChannelId:   filepath.Base(channelID),
		Name:        filepath.Base(name),
		Input:       input,
		maxFileSize: *a.srv.Config().FileSettings.MaxFileSize,
		maxImageRes: *a.srv.Config().FileSettings.MaxImageResolution,
		imgDecoder:  a.imgDecoder,
		imgEncoder:  a.imgEncoder,
	}
	// parse optional arguments:
	if userID != nil {
		t.UserId = *userID
	}
	if timestamp != nil {
		t.Timestamp = *timestamp
	}
	if contentLength != nil {
		t.ContentLength = *contentLength
	}
	if clientID != nil {
		t.ClientId = *clientID
	}
	if raw != nil {
		t.Raw = *raw
	}

	if *a.srv.Config().FileSettings.DriverName == "" {
		return nil, t.newAppError("api.file.upload_file.storage.app_error", http.StatusNotImplemented)
	}
	if t.ContentLength > t.maxFileSize {
		return nil, t.newAppError("api.file.upload_file.too_large_detailed.app_error", http.StatusRequestEntityTooLarge, "Length", t.ContentLength, "Limit", t.maxFileSize)
	}

	t.init(a)

	var aerr *model_helper.AppError
	if !t.Raw && model_helper.FileInfoIsImage(*t.fileinfo) {
		aerr = t.preprocessImage()
		if aerr != nil {
			return t.fileinfo, aerr
		}
	}

	written, aerr := t.writeFile(io.MultiReader(t.buf, t.limitedInput), t.fileinfo.Path)
	if aerr != nil {
		return nil, aerr
	}

	if written > t.maxFileSize {
		if fileErr := a.RemoveFile(t.fileinfo.Path); fileErr != nil {
			slog.Error("Failed to remove file", slog.Err(fileErr))
		}
		return nil, t.newAppError("api.file.upload_file.too_large_detailed.app_error", http.StatusRequestEntityTooLarge, "Length", t.ContentLength, "Limit", t.maxFileSize)
	}

	t.fileinfo.Size = written

	file, aerr := a.FileReader(t.fileinfo.Path)
	if aerr != nil {
		return nil, aerr
	}
	defer file.Close()

	// aerr = a.runPluginsHook(c, t.fileinfo, file)
	// if aerr != nil {
	// 	return nil, aerr
	// }

	if !t.Raw && model_helper.FileInfoIsImage(*t.fileinfo) {
		file, aerr = a.FileReader(t.fileinfo.Path)
		if aerr != nil {
			return nil, aerr
		}
		defer file.Close()
		t.postprocessImage(file)
	}

	if _, err := t.saveToDatabase(*t.fileinfo); err != nil {
		var appErr *model_helper.AppError
		switch {
		case errors.As(err, &appErr):
			return nil, appErr
		default:
			return nil, model_helper.NewAppError("UploadFileX", "app.file_info.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	if *a.srv.Config().FileSettings.ExtractContent {
		infoCopy := *t.fileinfo
		a.srv.Go(func() {
			err := a.ExtractContentFromFileInfo(&infoCopy)
			if err != nil {
				slog.Error("Failed to extract file content", slog.Err(err), slog.String("fileInfoId", infoCopy.ID))
			}
		})
	}

	return t.fileinfo, nil
}

func (t *UploadFileTask) preprocessImage() *model_helper.AppError {
	// If SVG, attempt to extract dimensions and then return
	if t.fileinfo.MimeType == "image/svg+xml" {
		svgInfo, err := imaging.ParseSVG(t.teeInput)
		if err != nil {
			slog.Warn("Failed to parse SVG", slog.Err(err))
		}
		if svgInfo.Width > 0 && svgInfo.Height > 0 {
			t.fileinfo.Width = model_types.NewNullInt(svgInfo.Width)
			t.fileinfo.Height = model_types.NewNullInt(svgInfo.Height)
		}
		t.fileinfo.HasPreviewImage = false
		return nil
	}

	// If we fail to decode, return "as is".
	w, h, err := imaging.GetDimensions(t.teeInput)
	if err != nil {
		return nil
	}
	t.fileinfo.Width = model_types.NewNullInt(w)
	t.fileinfo.Height = model_types.NewNullInt(h)

	if err = checkImageResolutionLimit(w, h, t.maxImageRes); err != nil {
		return t.newAppError("api.file.upload_file.large_image_detailed.app_error", http.StatusBadRequest)
	}

	t.fileinfo.HasPreviewImage = true
	nameWithoutExtension := t.Name[:strings.LastIndex(t.Name, ".")]
	t.fileinfo.PreviewPath = t.pathPrefix() + nameWithoutExtension + "_preview.jpg"
	t.fileinfo.ThumbnailPath = t.pathPrefix() + nameWithoutExtension + "_thumb.jpg"

	// check the image orientation with goexif; consume the bytes we
	// already have first, then keep Tee-ing from input.
	// TODO: try to reuse exif's .Raw buffer rather than Tee-ing
	if t.imageOrientation, err = imaging.GetImageOrientation(io.MultiReader(bytes.NewReader(t.buf.Bytes()), t.teeInput)); err == nil &&
		(t.imageOrientation == imaging.RotatedCWMirrored ||
			t.imageOrientation == imaging.RotatedCCW ||
			t.imageOrientation == imaging.RotatedCCWMirrored ||
			t.imageOrientation == imaging.RotatedCW) {
		t.fileinfo.Width, t.fileinfo.Height = t.fileinfo.Height, t.fileinfo.Width
	}

	// For animated GIFs disable the preview; since we have to Decode gifs
	// anyway, cache the decoded image for later.
	if t.fileinfo.MimeType == "image/gif" {
		gifConfig, err := gif.DecodeAll(io.MultiReader(bytes.NewReader(t.buf.Bytes()), t.teeInput))
		if err == nil {
			if len(gifConfig.Image) > 0 {
				t.fileinfo.HasPreviewImage = false
				t.decoded = gifConfig.Image[0]
				t.imageType = "gif"
			}
		}
	}

	return nil
}

func (t *UploadFileTask) postprocessImage(file io.Reader) {
	// don't try to process SVG files
	if t.fileinfo.MimeType == "image/svg+xml" {
		return
	}

	decoded, imgType := t.decoded, t.imageType
	if decoded == nil {
		var err error
		var release func()
		decoded, imgType, release, err = t.imgDecoder.DecodeMemBounded(file)
		if err != nil {
			slog.Error("Unable to decode image", slog.Err(err))
			return
		}
		defer release()
	}

	// Fill in the background of a potentially-transparent png file as white
	if imgType == "png" {
		imaging.FillImageTransparency(decoded, image.White)
	}

	decoded = imaging.MakeImageUpright(decoded, t.imageOrientation)
	if decoded == nil {
		return
	}

	writeJPEG := func(img image.Image, path string) {
		r, w := io.Pipe()
		go func() {
			err := t.imgEncoder.EncodeJPEG(w, img, jpegEncQuality)
			if err != nil {
				slog.Error("Unable to encode image as jpeg", slog.String("path", path), slog.Err(err))
				w.CloseWithError(err)
			} else {
				w.Close()
			}
		}()
		_, aerr := t.writeFile(r, path)
		if aerr != nil {
			slog.Error("Unable to upload", slog.String("path", path), slog.Err(aerr))
			return
		}
	}

	var wg sync.WaitGroup
	wg.Add(3)
	// Generating thumbnail and preview regardless of HasPreviewImage value.
	// This is needed on mobile in case of animated GIFs.
	go func() {
		defer wg.Done()
		writeJPEG(imaging.GenerateThumbnail(decoded, imageThumbnailWidth, imageThumbnailHeight), t.fileinfo.ThumbnailPath)
	}()

	go func() {
		defer wg.Done()
		writeJPEG(imaging.GeneratePreview(decoded, imagePreviewWidth), t.fileinfo.PreviewPath)
	}()

	go func() {
		defer wg.Done()
		if !t.fileinfo.MiniPreview.Valid || t.fileinfo.MiniPreview.Bytes == nil {
			if miniPreview, err := imaging.GenerateMiniPreviewImage(decoded, miniPreviewImageWidth, miniPreviewImageHeight, jpegEncQuality); err != nil {
				slog.Info("Unable to generate mini preview image", slog.Err(err))
			} else {
				t.fileinfo.MiniPreview = null.BytesFrom(miniPreview)
			}
		}
	}()
	wg.Wait()
}

func (t *UploadFileTask) pathPrefix() string {
	return t.Timestamp.Format("20060102") +
		// "/teams/" + t.TeamId +
		// "/channels/" + t.ChannelId +
		"/users/" + t.UserId +
		"/" + t.fileinfo.ID + "/"
}

func (t *UploadFileTask) newAppError(id string, httpStatus int, extra ...any) *model_helper.AppError {
	params := map[string]any{
		"Name":          t.Name,
		"Filename":      t.Name,
		"UserId":        t.UserId,
		"ContentLength": t.ContentLength,
		"ClientId":      t.ClientId,
		// "ChannelId":     t.ChannelId,
		// "TeamId":        t.TeamId,
	}
	if t.fileinfo != nil {
		params["Width"] = t.fileinfo.Width
		params["Height"] = t.fileinfo.Height
	}
	for i := 0; i+1 < len(extra); i += 2 {
		params[fmt.Sprintf("%v", extra[i])] = extra[i+1]
	}

	return model_helper.NewAppError("uploadFileTask", id, params, "", httpStatus)
}

func (a *ServiceFile) DoUploadFileExpectModification(c *request.Context, now time.Time, rawTeamId string, rawChannelId string, rawUserId string, rawFilename string, data []byte) (*model.FileInfo, []byte, *model_helper.AppError) {
	filename := filepath.Base(rawFilename)
	teamID := filepath.Base(rawTeamId)
	channelID := filepath.Base(rawChannelId)
	userID := filepath.Base(rawUserId)

	info, err := model_helper.GetInfoForBytes(filename, bytes.NewReader(data), len(data))
	if err != nil {
		err.StatusCode = http.StatusBadRequest
		return nil, data, err
	}

	if orientation, err := imaging.GetImageOrientation(bytes.NewReader(data)); err == nil &&
		(orientation == imaging.RotatedCWMirrored ||
			orientation == imaging.RotatedCCW ||
			orientation == imaging.RotatedCCWMirrored ||
			orientation == imaging.RotatedCW) {
		info.Width, info.Height = info.Height, info.Width
	}

	info.ID = model_helper.NewId()
	info.CreatorID = userID
	info.CreatedAt = now.UnixNano() / int64(time.Millisecond)

	pathPrefix := now.Format("20060102") + "/teams/" + teamID + "/channels/" + channelID + "/users/" + userID + "/" + info.ID + "/"
	info.Path = pathPrefix + filename

	if model_helper.FileInfoIsImage(*info) {
		if limitErr := checkImageResolutionLimit(model_helper.GetValueOfPointerOrZero(info.Width.Int), model_helper.GetValueOfPointerOrZero(info.Height.Int), *a.srv.Config().FileSettings.MaxImageResolution); limitErr != nil {
			err := model_helper.NewAppError("uploadFile", "api.file.upload_file.large_image.app_error", map[string]any{"Filename": filename}, limitErr.Error(), http.StatusBadRequest)
			return nil, data, err
		}

		nameWithoutExtension := filename[:strings.LastIndex(filename, ".")]
		info.PreviewPath = pathPrefix + nameWithoutExtension + "_preview.jpg"
		info.ThumbnailPath = pathPrefix + nameWithoutExtension + "_thumb.jpg"
	}

	if pluginsEnvironment, appErr := a.srv.Plugin.GetPluginsEnvironment(); appErr == nil && pluginsEnvironment != nil {
		var rejectionError *model_helper.AppError
		pluginContext := app.PluginContext(*c)
		pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
			var newBytes bytes.Buffer
			replacementInfo, rejectionReason := hooks.FileWillBeUploaded(pluginContext, info, bytes.NewReader(data), &newBytes)
			if rejectionReason != "" {
				rejectionError = model_helper.NewAppError("DoUploadFile", "File rejected by plugin. "+rejectionReason, nil, "", http.StatusBadRequest)
				return false
			}
			if replacementInfo != nil {
				info = replacementInfo
			}
			if newBytes.Len() != 0 {
				data = newBytes.Bytes()
				info.Size = int64(len(data))
			}

			return true
		}, plugin.FileWillBeUploadedID)
		if rejectionError != nil {
			return nil, data, rejectionError
		}
	}

	if _, err := a.WriteFile(bytes.NewReader(data), info.Path); err != nil {
		return nil, data, err
	}

	if _, err := a.srv.Store.FileInfo().Upsert(*info); err != nil {
		var appErr *model_helper.AppError
		switch {
		case errors.As(err, &appErr):
			return nil, data, appErr
		default:
			return nil, data, model_helper.NewAppError("DoUploadFileExpectModification", "app.file_info.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	if *a.srv.Config().FileSettings.ExtractContent {
		infoCopy := *info
		a.srv.Go(func() {
			err := a.ExtractContentFromFileInfo(&infoCopy)
			if err != nil {
				slog.Error("Failed to extract file content", slog.Err(err), slog.String("fileInfoId", infoCopy.ID))
			}
		})
	}

	return info, data, nil
}

func (a *ServiceFile) HandleImages(previewPathList []string, thumbnailPathList []string, fileData [][]byte) {
	var wg sync.WaitGroup

	for i := range fileData {
		img, release, err := prepareImage(a.imgDecoder, bytes.NewReader(fileData[i]))
		if err != nil {
			slog.Debug("Failed to prepare image", slog.Err(err))
			continue
		}
		wg.Add(2)
		go func(img image.Image, path string) {
			defer wg.Done()
			a.generateThumbnailImage(img, path)
		}(img, thumbnailPathList[i])

		go func(img image.Image, path string) {
			defer wg.Done()
			a.generatePreviewImage(img, path)
		}(img, previewPathList[i])

		wg.Wait()
		release()
	}
}

// prepareImage decodes raw fileData into an image, then returns that image, its width and height
func prepareImage(imgDecoder *imaging.Decoder, imgData io.ReadSeeker) (img image.Image, release func(), err error) {
	// Decode image bytes into Image object
	var imgType string
	img, imgType, release, err = imgDecoder.DecodeMemBounded(imgData)
	if err != nil {
		return nil, nil, fmt.Errorf("prepareImage: failed to decode image: %w", err)
	}

	// Fill in the background of a potentially-transparent png file as white
	if imgType == "png" {
		imaging.FillImageTransparency(img, image.White)
	}

	imgData.Seek(0, io.SeekStart)

	// Flip the image to be upright
	orientation, err := imaging.GetImageOrientation(imgData)
	if err != nil {
		slog.Debug("GetImageOrientation failed", slog.Err(err))
	}
	img = imaging.MakeImageUpright(img, orientation)

	return img, release, nil
}

func (a *ServiceFile) generateThumbnailImage(img image.Image, thumbnailPath string) {
	var buf bytes.Buffer
	if err := a.imgEncoder.EncodeJPEG(&buf, imaging.GenerateThumbnail(img, imageThumbnailWidth, imageThumbnailHeight), jpegEncQuality); err != nil {
		slog.Error("Unable to encode image as jpeg", slog.String("path", thumbnailPath), slog.Err(err))
		return
	}

	if _, err := a.WriteFile(&buf, thumbnailPath); err != nil {
		slog.Error("Unable to upload thumbnail", slog.String("path", thumbnailPath), slog.Err(err))
		return
	}
}

func (a *ServiceFile) generatePreviewImage(img image.Image, previewPath string) {
	var buf bytes.Buffer
	preview := imaging.GeneratePreview(img, imagePreviewWidth)

	if err := a.imgEncoder.EncodeJPEG(&buf, preview, jpegEncQuality); err != nil {
		slog.Error("Unable to encode image as preview jpg", slog.Err(err), slog.String("path", previewPath))
		return
	}

	if _, err := a.WriteFile(&buf, previewPath); err != nil {
		slog.Error("Unable to upload preview", slog.Err(err), slog.String("path", previewPath))
		return
	}
}

// generateMiniPreview updates mini preview if needed
// will save fileinfo with the preview added
func (a *ServiceFile) generateMiniPreview(fi *model.FileInfo) {
	if model_helper.FileInfoIsImage(*fi) && !fi.MiniPreview.Valid {
		file, appErr := a.FileReader(fi.Path)
		if appErr != nil {
			slog.Debug("error reading image file", slog.Err(appErr))
			return
		}
		defer file.Close()
		img, release, err := prepareImage(a.imgDecoder, file)
		if err != nil {
			slog.Debug(
				"generateMiniPreview: prepareImage failed",
				slog.Err(err),
				slog.String("fileinfo_id", fi.ID),
				slog.String("creator_id", fi.CreatorID),
			)
			// Since this file is not a valid image (for whatever reason), prevent this fileInfo
			// from entering generateMiniPreview in the future
			fi.UpdatedAt = model_helper.GetMillis()
			fi.MimeType = "invalid-" + fi.MimeType
			if _, err = a.srv.Store.FileInfo().Upsert(*fi); err != nil {
				slog.Debug("Invalidating FileInfo failed", slog.Err(err))
			}
			return
		}
		defer release()
		var miniPreview []byte
		if miniPreview, err = imaging.GenerateMiniPreviewImage(
			img,
			miniPreviewImageWidth,
			miniPreviewImageHeight,
			jpegEncQuality,
		); err != nil {
			slog.Info("Unable to generate mini preview image", slog.Err(err))
		} else {
			fi.MiniPreview = null.BytesFrom(miniPreview)
		}
		if _, err = a.srv.Store.FileInfo().Upsert(*fi); err != nil {
			slog.Debug("creating mini preview failed", slog.Err(err))
		} else {
			// a.srv.Store.FileInfo().InvalidateFileInfosForPostCache(fi.PostId, false)
			slog.Warn("this flow is not implemented", slog.String("function", "generateMiniPreview"))
		}
	}
}

func (a *ServiceFile) generateMiniPreviewForInfos(fileInfos []*model.FileInfo) {
	wg := new(sync.WaitGroup)
	wg.Add(len(fileInfos))

	for _, fileInfo := range fileInfos {
		go func(fi *model.FileInfo) {
			defer wg.Done()
			a.generateMiniPreview(fi)
		}(fileInfo)
	}
	wg.Wait()
}

// GetFileInfo get fileInfo object from database with given fileID, populates its "MiniPreview" and returns it.
func (a *ServiceFile) GetFileInfos(page, perPage int, opt model_helper.FileInfoFilterOption) ([]*model.FileInfo, *model_helper.AppError) {
	fileInfos, err := a.srv.Store.FileInfo().GetWithOptions(opt)
	if err != nil {
		var invErr *store.ErrInvalidInput
		var ltErr *store.ErrLimitExceeded
		switch {
		case errors.As(err, &invErr):
			return nil, model_helper.NewAppError("GetFileInfos", "app.file_info.get_with_options.app_error", nil, invErr.Error(), http.StatusBadRequest)
		case errors.As(err, &ltErr):
			return nil, model_helper.NewAppError("GetFileInfos", "app.file_info.get_with_options.app_error", nil, ltErr.Error(), http.StatusBadRequest)
		default:
			return nil, model_helper.NewAppError("GetFileInfos", "app.file_info.get_with_options.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	a.generateMiniPreviewForInfos(fileInfos)

	return fileInfos, nil
}

func (a *ServiceFile) GetFileInfo(fileID string) (*model.FileInfo, *model_helper.AppError) {
	fileInfo, err := a.srv.Store.FileInfo().Get(fileID, false)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model_helper.NewAppError("GetFileInfo", "app.file_info.get.app_error", nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model_helper.NewAppError("GetFileInfo", "app.file_info.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	a.generateMiniPreview(fileInfo)
	return fileInfo, nil
}

func (a *ServiceFile) GetFile(fileID string) ([]byte, *model_helper.AppError) {
	info, err := a.GetFileInfo(fileID)
	if err != nil {
		return nil, err
	}

	data, err := a.ReadFile(info.Path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (a *ServiceFile) CopyFileInfos(userID string, fileIDs []string) ([]string, *model_helper.AppError) {
	var newFileIds []string

	now := model_helper.GetMillis()

	for _, fileID := range fileIDs {
		fileInfo, err := a.srv.Store.FileInfo().Get(fileID, false)
		if err != nil {
			var nfErr *store.ErrNotFound
			switch {
			case errors.As(err, &nfErr):
				return nil, model_helper.NewAppError("CopyFileInfos", "app.file_info.get.app_error", nil, nfErr.Error(), http.StatusNotFound)
			default:
				return nil, model_helper.NewAppError("CopyFileInfos", "app.file_info.get.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}

		fileInfo.ID = model_helper.NewId()
		fileInfo.CreatorID = userID
		fileInfo.CreatedAt = now
		fileInfo.UpdatedAt = now

		if _, err := a.srv.Store.FileInfo().Upsert(*fileInfo); err != nil {
			var appErr *model_helper.AppError
			switch {
			case errors.As(err, &appErr):
				return nil, appErr
			default:
				return nil, model_helper.NewAppError("CopyFileInfos", "app.file_info.save.app_error", nil, err.Error(), http.StatusInternalServerError)
			}
		}

		newFileIds = append(newFileIds, fileInfo.ID)
	}

	return newFileIds, nil
}

// This function zip's up all the files in fileDatas array and then saves it to the directory specified with the specified zip file name
// Ensure the zip file name ends with a .zip
func (a *ServiceFile) CreateZipFileAndAddFiles(fileBackend filestore.FileBackend, fileDatas []model_helper.FileData, zipFileName, directory string) error {
	// Create Zip File (temporarily stored on disk)
	conglomerateZipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer os.Remove(zipFileName)

	// Create a new zip archive.
	zipFileWriter := zip.NewWriter(conglomerateZipFile)

	// Populate Zip file with File Datas array
	err = populateZipfile(zipFileWriter, fileDatas)
	if err != nil {
		return err
	}

	conglomerateZipFile.Seek(0, 0)
	_, err = fileBackend.WriteFile(conglomerateZipFile, path.Join(directory, zipFileName))
	if err != nil {
		return err
	}

	return nil
}

// This is a implementation of Go's example of writing files to zip (with slight modification)
// https://golang.org/src/archive/zip/example_test.go
func populateZipfile(w *zip.Writer, fileDatas []model_helper.FileData) error {
	defer w.Close()
	for _, fd := range fileDatas {
		f, err := w.Create(fd.Filename)
		if err != nil {
			return err
		}

		_, err = f.Write(fd.Body)
		if err != nil {
			return err
		}
	}
	return nil
}
