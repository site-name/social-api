package markup

// import (
// 	"regexp"
// 	"strings"
// 	"unicode"
// 	"unicode/utf8"
// )

// var (
// 	commentClose = regexp.MustCompile(`--\s*>`)
// )

// type ParserBase struct {
// 	lineno  int64
// 	offset  int64
// 	rawData string
// }

// func (p *ParserBase) Error(message string) {
// 	panic(message)
// }

// func (p *ParserBase) Reset() {
// 	p.lineno = 1
// 	p.offset = 0
// }

// // Return current line number and offset.
// func (p *ParserBase) GetPos() (lineNumber int64, offet int64) {
// 	return p.lineno, p.offset
// }

// func (p *ParserBase) UpdatePos(i, j int64) int64 {
// 	if i >= j {
// 		return j
// 	}

// 	nlines := strings.Count(p.rawData[i:j], "\n")
// 	if nlines > 0 {
// 		p.lineno += int64(nlines)
// 		pos := strings.LastIndex(p.rawData[i:j], "\n")
// 		p.offset = j - (int64(pos) + 1)
// 	} else {
// 		p.offset += (j - i)
// 	}
// 	return j
// }

// func (p *ParserBase) ParseDeclaration(i int64) int64 {

// 	// This is some sort of declaration; in "HTML as
// 	// deployed," this should only be the document type
// 	// declaration ("<!DOCTYPE html...>").
// 	// ISO 8879:1986, however, has more complex
// 	// declaration syntax for elements in <!...>, including:
// 	// --comment--
// 	// [marked section]
// 	// name in the following list: ENTITY, DOCTYPE, ELEMENT,
// 	// ATTLIST, NOTATION, SHORTREF, USEMAP,
// 	// LINKTYPE, LINK, IDLINK, USELINK, SYSTEM

// 	j := i + 2
// 	if p.rawData[i:j] != "<!" {
// 		panic("unexpected call to parse_declaration")
// 	}
// 	if p.rawData[j:j+1] == ">" {
// 		// the empty comment <!>
// 		return j + 1
// 	}

// 	if strings.ContainsAny(p.rawData[j:j+1], "-" + "") {
// 		// Start of comment followed by buffer boundary,
//     // or just a buffer boundary.
// 		return -1
// 	}
// 	n := len(p.rawData)
// 	if p.rawData[j:j+2] == "--" {
// 		// comment
// 		return p.ParseComment(i)
// 	}
// }

// func (p *ParserBase) ParseComment(i int64, report bool) int64 {
// 	if p.rawData[i:i+4] != "<!--" {
// 		p.Error("unexpected call to ParseComment()")
// 	}
// 	commentClose.FindAllString(p.rawData)
// }
