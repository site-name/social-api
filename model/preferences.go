package model

import (
	"io"

	"github.com/sitename/sitename/modules/json"
)

type Preferences []Preference

func (o *Preferences) ToJson() string {
	b, _ := json.JSON.Marshal(o)
	return string(b)
}

func PreferencesFromJson(data io.Reader) (Preferences, error) {
	decoder := json.JSON.NewDecoder(data)
	var o Preferences
	err := decoder.Decode(&o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
