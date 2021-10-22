package file

import (
	"io"

	"github.com/sitename/sitename/model"
)

type FileInfoSearchMatches map[string][]string

type FileInfoSearchResults struct {
	*FileInfoList
	Matches FileInfoSearchMatches `json:"matches"`
}

func MakeFileInfoSearchResults(fileInfos *FileInfoList, matches FileInfoSearchMatches) *FileInfoSearchResults {
	return &FileInfoSearchResults{
		fileInfos,
		matches,
	}
}

func (o *FileInfoSearchResults) ToJSON() string {
	return model.ModelToJson(o)
}

func FileInfoSearchResultsFromJson(data io.Reader) *FileInfoSearchResults {
	var o *FileInfoSearchResults
	model.ModelFromJson(&o, data)
	return o
}
