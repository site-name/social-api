package store

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/modules/slog"
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

type Rollbackable interface {
	Rollback() error
}

// finalizeTransaction ensures a transaction is closed after use, rolling back if not already committed.
func FinalizeTransaction(transaction Rollbackable) {
	if err := transaction.Rollback(); err != nil {
		slog.Error("Failed to rollback transaction", slog.Err(err))
	}
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
// 		"quoted string"  -  will return true
// 		unquoted string  -  will return false
func IsQuotedWord(s string) bool {
	if len(s) < 2 {
		return false
	}

	return s[0] == '"' && s[len(s)-1] == '"'
}

// WildcardSearchTerm convert given term to lower-case, concatenates `%` to both ends
//
// Example:
//  WildcardSearchTerm("HELLO") => "%hello%"
func WildcardSearchTerm(term string) string {
	return strings.ToLower("%" + term + "%")
}

// AppErrorFromDatabaseLookupError is a utility function that create *model.AppError with given error.
//
// Must be used with database LOOLUP errors.
func AppErrorFromDatabaseLookupError(where, errId string, err error) *model.AppError {
	statusCode := http.StatusInternalServerError
	var nfErr *ErrNotFound
	if errors.As(err, &nfErr) {
		statusCode = http.StatusNotFound
	}

	return model.NewAppError(where, errId, nil, err.Error(), statusCode)
}
