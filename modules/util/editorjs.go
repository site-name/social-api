package util

import (
	"regexp"
	// "strings"
)

var (
	BLACKLISTED_URL_SCHEMES        = []string{"javascript"}
	HYPERLINK_TAG_WITH_URL_PATTERN = regexp.MustCompile(`(.*?<a\s+href=\\?\")(\w+://\S+[^\\])(\\?\">)`)
)

// Sanitize a given EditorJS JSON definitions.
//
// Look for not allowed URLs, replaced them with `invalid` value, and clean valid ones.
//
// `to_string` flag is used for returning concatenated string from all blocks
// instead of returning json object.
// func CleanEditorJS(definitions map[string][]map[string]interface{}) string {
// 	var str string

// 	blocks, found := definitions["blocks"]
// 	if !found {
// 		return str
// 	}

// 	for key, block := range blocks {
// 		blockType := block["type"]
// 		data, ok := block["data"]
// 		dataMap, yes := data.(map[string]interface{})
// 		if !ok || data == nil || !yes {
// 			continue
// 		}

// 		if blockType == "list" {
// 			items := dataMap["items"].([]interface{})
// 			for itemIdx, item := range items {
// 				if item == nil {
// 					continue
// 				}

// 				strItem := item.(string)
// 				newText := CleanTextData(strItem)
// 				// string +=
// 			}
// 		} else {

// 		}
// 	}
// }

// Look for url in text, check if URL is allowed and return the cleaned URL.
//
// By default, only the protocol ``javascript`` is denied.
// func CleanTextData(text string) string {
// 	if text == "" {
// 		return ""
// 	}

// 	var endOfMatch int
// 	var newText string

// 	matches := HYPERLINK_TAG_WITH_URL_PATTERN.FindAllString(text, -1)

// }

// Internal tag stripping utility used by strip_tags.
// func stripOnce(value string) string {

// }

// Return the given HTML with all tags stripped.
// func StripTags(value string) string {
// 	for strings.ContainsRune(value, '<') && strings.ContainsRune(value, '>') {
// 		newValue := stripOnce(value)
// 	}
// }
