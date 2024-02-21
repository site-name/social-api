package docextractor

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/pkg/errors"
)

type snPreviewExtractor struct {
	url          string
	secret       string
	pdfExtractor pdfExtractor
}

var snPreviewSupportedExtensions = map[string]bool{
	"ppt":  true,
	"odp":  true,
	"xls":  true,
	"xlsx": true,
	"ods":  true,
}

func newSnPreviewExtractor(url string, secret string, pdfExtractor pdfExtractor) *snPreviewExtractor {
	return &snPreviewExtractor{url, secret, pdfExtractor}
}

func (sn *snPreviewExtractor) Match(filename string) bool {
	extension := strings.TrimPrefix(path.Ext(filename), ".")
	return snPreviewSupportedExtensions[extension]
}

func (sn *snPreviewExtractor) Extract(filename string, file io.ReadSeeker) (string, error) {
	b, w, err := createMultipartFormDate("file", filename, file)
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate file preview using snPreview.")
	}
	req, err := http.NewRequest("POST", sn.url+"/toPDF", &b)
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate file preview using snPreview.")
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if sn.secret != "" {
		req.Header.Add("Authentication", sn.secret)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate file preview using mmpreview.")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Unable to generate file preview using mmpreview (The server has replied with an error)")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read the response from snPreview")
	}
	return sn.pdfExtractor.Extract(filename, bytes.NewReader(data))
}

func createMultipartFormDate(fieldName, fileName string, fileData io.ReadSeeker) (bytes.Buffer, *multipart.Writer, error) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	if fw, err = w.CreateFormFile(fieldName, fileName); err != nil {
		return b, nil, err
	}
	if _, err = io.Copy(fw, fileData); err != nil {
		return b, nil, err
	}
	w.Close()
	return b, w, nil
}
