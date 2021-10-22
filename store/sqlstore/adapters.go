package sqlstore

import (
	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/mattermost/gorp"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/account"
	"github.com/sitename/sitename/modules/i18n"
	"github.com/sitename/sitename/modules/json"
)

// siteNameConverter make tables able to have fields with custom types
//
// Example:
//  map[string]string, []string, map[string]interface{}, ...
type siteNameConverter struct{}

func (me siteNameConverter) ToDb(val interface{}) (interface{}, error) {
	switch t := val.(type) {
	case model.StringMap:
		return model.MapToJson(t), nil
	case account.StringMap: // this is needed
		return model.MapToJson(t), nil
	case map[string]string:
		return model.MapToJson(t), nil
	case model.StringArray:
		return model.ArrayToJson(t), nil
	case model.StringInterface:
		return model.StringInterfaceToJson(t), nil
	case map[string]interface{}:
		return model.StringInterfaceToJson(t), nil
	case JSONSerializable:
		return t.ToJSON(), nil
	case *opengraph.OpenGraph:
		return json.JSON.Marshal(t)
	}

	return val, nil
}

func (me siteNameConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	switch target.(type) {
	case *model.StringMap, *account.StringMap, *map[string]string:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(i18n.T("store.sql.convert_string_map"))
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.StringArray, *[]string:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(i18n.T("store.sql.convert_string_array"))
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	case *model.StringInterface, *map[string]interface{}:
		binder := func(holder, target interface{}) error {
			s, ok := holder.(*string)
			if !ok {
				return errors.New(i18n.T("store.sql.convert_string_interface"))
			}
			b := []byte(*s)
			return json.JSON.Unmarshal(b, target)
		}
		return gorp.CustomScanner{Holder: new(string), Target: target, Binder: binder}, true
	}

	return gorp.CustomScanner{}, false
}

type JSONSerializable interface {
	ToJSON() string
}
