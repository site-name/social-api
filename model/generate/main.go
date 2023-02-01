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

	"github.com/gosimple/slug"
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
	Id             string         `json:"id"`
	VnName         string         `json:"vn_name"`
	EnName         string         `json:"en_name"`
	Slug           string         `json:"slug"`
	Images         []string       `json:"images"`
	SeoTitle       *string        `json:"seo_title"`
	SeoDescription *string        `json:"seo_description"`
	Description    map[string]any `json:"description"`
	Children       []string       `json:"children"`

	Named string `json:"named"`
}

var firstLevel = map[string]map[string]struct{}{}
var secondLevel = map[string]map[string]struct{}{}
var thirdLevel = map[string]map[string]struct{}{}
var fourthLevel = map[string]map[string]struct{}{}

type data struct {
	Categories  []*DesiredCategory
	FirstLevels []string
}

var replacer = strings.NewReplacer(
	",", "",
	"&", "",
	"-", "",
	"/", "",
	"'", "",
	" ", "",
	">", "",
	"+", "plus",
)

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

	meetMap := map[string]struct{}{}

	desireds := []*DesiredCategory{}
	for idx, cate := range cates {
		if cate.CategoryNameEn == "" {
			continue
		}

		var firstLevelName, secondLevelName,
			thirdLevelName, fourthLevelName,
			fifthLevelName string

		named := "Category"
		slugg := ""
		for pathIdx, path := range cate.Path {

			named += replacer.Replace(path.CategoryNameEn)
			slugg += " " + path.CategoryNameEn

			switch pathIdx {
			case 0:
				firstLevelName = named
			case 1:
				secondLevelName = named
			case 2:
				thirdLevelName = named
			case 3:
				fourthLevelName = named
			case 4:
				fifthLevelName = named
			}

			if pathIdx == len(cate.Path)-1 {
				if secondLevelName != "" {
					if firstLevel[firstLevelName] == nil {
						firstLevel[firstLevelName] = map[string]struct{}{}
					}
					firstLevel[firstLevelName][secondLevelName] = struct{}{}
				}
				if thirdLevelName != "" {
					if secondLevel[secondLevelName] == nil {
						secondLevel[secondLevelName] = map[string]struct{}{}
					}
					secondLevel[secondLevelName][thirdLevelName] = struct{}{}
				}
				if fourthLevelName != "" {
					if thirdLevel[thirdLevelName] == nil {
						thirdLevel[thirdLevelName] = map[string]struct{}{}
					}
					thirdLevel[thirdLevelName][fourthLevelName] = struct{}{}
				}
				if fifthLevelName != "" {
					if fourthLevel[fourthLevelName] == nil {
						fourthLevel[fourthLevelName] = map[string]struct{}{}
					}
					fourthLevel[fourthLevelName][fifthLevelName] = struct{}{}
				}
			}

			if _, met := meetMap[named]; !met {
				desired := &DesiredCategory{
					Id:     strconv.Itoa(idx) + "x" + strconv.Itoa(pathIdx),
					Named:  named,
					Slug:   slug.Make(slugg),
					VnName: path.CategoryName,
					EnName: path.CategoryNameEn,
				}
				if pathIdx == len(cate.Path)-1 {
					desired.Images = cate.Images
				}

				desireds = append(desireds, desired)
				meetMap[named] = struct{}{}
			}
		}
	}

	firstLevels := []string{}

	for _, cate := range desireds {
		switch {
		case firstLevel[cate.Named] != nil:
			for key := range firstLevel[cate.Named] {
				cate.Children = append(cate.Children, key)
			}
			firstLevels = append(firstLevels, cate.Named)

		case secondLevel[cate.Named] != nil:
			for key := range secondLevel[cate.Named] {
				cate.Children = append(cate.Children, key)
			}
		case thirdLevel[cate.Named] != nil:
			for key := range thirdLevel[cate.Named] {
				cate.Children = append(cate.Children, key)
			}
		case fourthLevel[cate.Named] != nil:
			for key := range fourthLevel[cate.Named] {
				cate.Children = append(cate.Children, key)
			}
		}
	}

	out := bytes.NewBufferString("")

	t := template.Must(template.New("t.go.tmpl").ParseFiles("t.go.tmpl"))
	if err = t.Execute(out, data{Categories: desireds, FirstLevels: firstLevels}); err != nil {
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
