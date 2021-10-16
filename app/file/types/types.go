package types

import (
	"bytes"
	"image"
	"io"
	"time"

	"github.com/sitename/sitename/app/imaging"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/modules/plugin"
)

type UploadFileTask struct {
	Name               string
	UserId             string
	Timestamp          time.Time      // Time stamp to use when creating the file.
	ContentLength      int64          // The value of the Content-Length http header, when available.
	Input              io.Reader      // The file data stream.
	ClientId           string         // An optional, client-assigned Id field.
	Raw                bool           // If Raw, do not execute special processing for images, just upload the file.  Plugins are still invoked.
	buf                *bytes.Buffer  //
	limit              int64          //
	limitedInput       io.Reader      //
	teeInput           io.Reader      //
	fileinfo           *file.FileInfo //
	maxFileSize        int64          //
	maxImageRes        int64          //
	decoded            image.Image    // Cached image data that (may) get initialized in preprocessImage and is used in postprocessImage
	imageType          string
	imageOrientation   int
	writeFile          func(io.Reader, string) (int64, *model.AppError)
	saveToDatabase     func(*file.FileInfo) (*file.FileInfo, error)
	imgDecoder         *imaging.Decoder
	imgEncoder         *imaging.Encoder
	pluginsEnvironment *plugin.Environment

	// Testing: overrideable dependency functions
	// ChannelId string
	// TeamId    string
}
