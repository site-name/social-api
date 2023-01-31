//go:generate go run main.go

package main

import (
	"bytes"
	"encoding/json"
	"go/format"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

type CategoryPath struct {
	CategoryID     int    `json:"category_id"`
	CategoryName   string `json:"category_name"`
	CategoryNameEn string `json:"category_name_en"`
}

type Category struct {
	CategoryID     int            `json:"category_id"`
	CategoryName   string         `json:"category_name"`
	CategoryNameEn string         `json:"category_name_en"`
	Toggle         bool           `json:"toggle"`
	Images         []string       `json:"images"`
	Path           []CategoryPath `json:"path"`
}

type DesiredCategory struct {
	Id     string `json:"id"`
	VnName string `json:"vn_name"`
	// E.g:
	//  {"en": "Women Clothes", "vi": "Thời Trang Nữ"}
	Name           map[string]string  `json:"name"`
	Slug           string             `json:"slug"`
	Images         []string           `json:"images"`
	SeoTitle       *string            `json:"seo_title"`
	SeoDescription *string            `json:"seo_description"`
	Description    map[string]any     `json:"description"`
	Children       []*DesiredCategory `json:"children"`

	Named string `json:"named"`
}

// var firstLevel = map[string][]string{}
// var secondLevel = map[string][]string{}
// var thirdLevel = map[string][]string{}

type data struct {
	Categories []*DesiredCategory
}

func main() {
	rawData, err := os.ReadFile("./raw_categories.json")
	if err != nil {
		log.Fatalln("Error reading json file: ", err)
	}

	var cates []*Category
	err = json.Unmarshal(rawData, &cates)
	if err != nil {
		log.Fatalln("Error unmarshaling:", err)
	}

	desireds := []*DesiredCategory{}
	for idx, cate := range cates {
		if cate.CategoryNameEn == "" {
			continue
		}

		named := "Category"
		slug := ""
		for _, path := range cate.Path {
			split := strings.FieldsFunc(path.CategoryNameEn, func(r rune) bool {
				return r == ',' || r == '&' || unicode.IsSpace(r) || r == '-' || r == '/'
			})

			for _, item := range split {
				item = strings.ReplaceAll(item, "'", "")
				named += item

				lowerItem := strings.ToLower(item)
				if slug != "" {
					lowerItem = "-" + lowerItem
				}
				slug += lowerItem
			}
		}

		d := &DesiredCategory{
			Id: strconv.Itoa(idx + 1),
			Name: map[string]string{
				"en": cate.CategoryNameEn,
				"vi": cate.CategoryName,
			},
			VnName: cate.CategoryName,
			Images: cate.Images,
			Slug:   slug,
			Named:  named,
		}
		desireds = append(desireds, d)
	}

	out := bytes.NewBufferString("")

	t := template.Must(template.New("t.go.tmpl").ParseFiles("t.go.tmpl"))
	if err = t.Execute(out, data{Categories: desireds}); err != nil {
		log.Fatalln("error template:", err)
	}

	source, err := format.Source(out.Bytes())
	if err != nil {
		log.Fatalln("error formatting source code:", err)
	}

	err = os.WriteFile("../CATEGORIES.go", source, 0644)
	if err != nil {
		log.Fatalln("error writing file:", err)
	}
}
