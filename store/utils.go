package store

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/sitename/sitename/modules/util"
)

var escapeLikeSearchChar = []string{
	"%",
	"_",
}

func SanitizeSearchTerm(term string, escapeChar string) string {
	term = strings.Replace(term, escapeChar, "", -1)

	for _, c := range escapeLikeSearchChar {
		term = strings.Replace(term, c, escapeChar+c, -1)
	}

	return term
}

// Converts a list of strings into a list of query parameters and a named parameter map that can
// be used as part of a SQL query.
func MapStringsToQueryParams(list []string, paramPrefix string) (string, map[string]interface{}) {
	var keys strings.Builder
	params := make(map[string]interface{}, len(list))
	for i, entry := range list {
		if keys.Len() > 0 {
			keys.WriteString(",")
		}

		key := paramPrefix + strconv.Itoa(i)
		keys.WriteString(":" + key)
		params[key] = entry
	}

	return "(" + keys.String() + ")", params
}

// removeNonAlphaNumericUnquotedTerms removes all unquoted words that only contain
// non-alphanumeric chars from given line
func RemoveNonAlphaNumericUnquotedTerms(line, separator string) string {
	words := strings.Split(line, separator)
	filteredResult := make([]string, 0, len(words))

	for _, w := range words {
		if IsQuotedWord(w) || ContainsAlphaNumericChar(w) {
			filteredResult = append(filteredResult, strings.TrimSpace(w))
		}
	}
	return strings.Join(filteredResult, separator)
}

// containsAlphaNumericChar returns true in case any letter or digit is present, false otherwise
func ContainsAlphaNumericChar(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// isQuotedWord return true if the input string is quoted, false otherwise. Ex :-
//
//	"quoted string"  -  will return true
//	unquoted string  -  will return false
func IsQuotedWord(s string) bool {
	if len(s) < 2 {
		return false
	}

	return s[0] == '"' && s[len(s)-1] == '"'
}

// WildcardSearchTerm convert given term to lower-case, concatenates `%` to both ends
//
// Example:
//
//	WildcardSearchTerm("HELLO") => "%hello%"
func WildcardSearchTerm(term string) string {
	return strings.ToLower("%" + term + "%")
}

// SqlizerIsEqualNull checks if given expr is like squirrel.Eq{"...": nil}
func SqlizerIsEqualNull(expr squirrel.Sqlizer) bool {
	eq, ok := expr.(squirrel.Eq)
	if ok {
		for _, value := range eq {
			if value == nil {
				return true
			}
		}
		return false
	}

	return false
}

// Eg:
//
//	type Extra struct {
//	  Embed string
//	}
//	type Person struct {
//	  Name string
//	  Private string `db:"-"`
//	  unExported int
//	  Extra
//	}
//
//	p := &Person{}
//	ExtractModelFieldPointers(p) == []any{&p.Name, &p.Embed}
func ExtractModelFieldPointers(modelPointer any) []any {
	valueOf := reflect.ValueOf(modelPointer)

	if valueOf.Kind() != reflect.Pointer || valueOf.Elem().Kind() != reflect.Struct {
		panic("obj must be a pointer to a struct model")
	}

	res := []any{}
	for _, fieldName := range ExtractModelFieldNames(valueOf.Elem().Interface()) {
		fieldPointer := valueOf.Elem().FieldByName(fieldName).Addr().Interface()
		res = append(res, fieldPointer)
	}

	return res
}

// Eg:
//
//	type Extra struct {
//	  Embed string
//	}
//	type Person struct {
//	  Name string
//	  Private string `db:"-"`
//	  unExported int
//	  Extra
//	}
//
//	ExtractModelFieldNames(Person{}) == []string{"Name", "Embed"}
func ExtractModelFieldNames(model any) util.AnyArray[string] {
	res := []string{}
	valueOf := reflect.ValueOf(model)
	typeOf := reflect.TypeOf(model)

	for i := 0; i < valueOf.NumField(); i++ {
		fieldTypeAtIdx := typeOf.Field(i)

		switch {
		case !fieldTypeAtIdx.IsExported() ||
			fieldTypeAtIdx.Tag.Get("db") == "-":
			continue

		case fieldTypeAtIdx.Type.Kind() == reflect.Struct:
			fieldValueAtIdx := valueOf.Field(i)
			res = append(res, ExtractModelFieldNames(fieldValueAtIdx.Interface())...)

		default:
			res = append(res, fieldTypeAtIdx.Name)
		}
	}

	return res
}

func BuildSqlizer(option squirrel.Sqlizer, where string) ([]any, error) {
	if option == nil {
		return []any{}, nil
	}

	query, args, err := option.ToSql()
	if err != nil {
		return []any{}, errors.Wrap(err, where+"_ToSql")
	}

	res := make([]any, 0, len(args)+1)
	res[0] = query
	return append(res, args...), nil
}
