package sqlstore

import (
	"database/sql"
	"strconv"
	"strings"
	"unicode"

	"github.com/mattermost/gorp"
	"github.com/sitename/sitename/modules/log"
)

var escapeLikeSearchChar = []string{
	"%",
	"_",
}

func sanitizeSearchTerm(term string, escapeChar string) string {
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

// finalizeTransaction ensures a transaction is closed after use, rolling back if not already committed.
func finalizeTransaction(transaction interface{}) {
	// Rollback returns sql.ErrTxDone if the transaction was already closed.
	switch t := transaction.(type) {
	case *gorp.Transaction:
		if err := t.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Error("Failed to rollback transaction: %v", err)
		}
	case *sql.Tx:
		if err := t.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Error("Failed to rollback transaction: %v", err)
		}
	}
}

// removeNonAlphaNumericUnquotedTerms removes all unquoted words that only contain
// non-alphanumeric chars from given line
func removeNonAlphaNumericUnquotedTerms(line, separator string) string {
	words := strings.Split(line, separator)
	filteredResult := make([]string, 0, len(words))

	for _, w := range words {
		if isQuotedWord(w) || containsAlphaNumericChar(w) {
			filteredResult = append(filteredResult, strings.TrimSpace(w))
		}
	}
	return strings.Join(filteredResult, separator)
}

// containsAlphaNumericChar returns true in case any letter or digit is present, false otherwise
func containsAlphaNumericChar(s string) bool {
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
func isQuotedWord(s string) bool {
	if len(s) < 2 {
		return false
	}

	return s[0] == '"' && s[len(s)-1] == '"'
}

func wildcardSearchTerm(term string) string {
	return strings.ToLower("%" + term + "%")
}
