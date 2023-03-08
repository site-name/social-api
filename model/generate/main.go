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
	Categories   []*DesiredCategory
	FirstLevels  []string
	SecondLevels []string
	ThirdLevels  []string
	FourthLevels []string
	FifthLevels  []string
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
	dt := data{}

	for idx, cate := range cates {
		if cate.CategoryNameEn == "" {
			continue
		}

		var (
			levelNameMap = map[int]string{}
			named        = "Category"
			slugg        = ""
		)
		for pathIdx, path := range cate.Path {

			named += replacer.Replace(path.CategoryNameEn)
			slugg += " " + path.CategoryNameEn
			levelNameMap[pathIdx] = named

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

				dt.Categories = append(dt.Categories, desired)
				switch pathIdx {
				case 0:
					dt.FirstLevels = append(dt.FirstLevels, named)
				case 1:
					dt.SecondLevels = append(dt.SecondLevels, named)
				case 2:
					dt.ThirdLevels = append(dt.ThirdLevels, named)
				case 3:
					dt.FourthLevels = append(dt.FourthLevels, named)
				case 4:
					dt.FifthLevels = append(dt.FifthLevels, named)
				}
				meetMap[named] = struct{}{}
			}
		}

		if levelNameMap[1] != "" {
			if firstLevel[levelNameMap[0]] == nil {
				firstLevel[levelNameMap[0]] = map[string]struct{}{}
			}
			firstLevel[levelNameMap[0]][levelNameMap[1]] = struct{}{}
		}
		if levelNameMap[2] != "" {
			if secondLevel[levelNameMap[1]] == nil {
				secondLevel[levelNameMap[1]] = map[string]struct{}{}
			}
			secondLevel[levelNameMap[1]][levelNameMap[2]] = struct{}{}
		}
		if levelNameMap[3] != "" {
			if thirdLevel[levelNameMap[2]] == nil {
				thirdLevel[levelNameMap[2]] = map[string]struct{}{}
			}
			thirdLevel[levelNameMap[2]][levelNameMap[3]] = struct{}{}
		}
		if levelNameMap[4] != "" {
			if fourthLevel[levelNameMap[3]] == nil {
				fourthLevel[levelNameMap[3]] = map[string]struct{}{}
			}
			fourthLevel[levelNameMap[3]][levelNameMap[4]] = struct{}{}
		}
	}

	for _, cate := range dt.Categories {
		switch {
		case firstLevel[cate.Named] != nil:
			for key := range firstLevel[cate.Named] {
				cate.Children = append(cate.Children, key)
			}
			// dt.FirstLevels = append(dt.FirstLevels, cate.Named)

		case secondLevel[cate.Named] != nil:
			for key := range secondLevel[cate.Named] {
				cate.Children = append(cate.Children, key)
			}
			// dt.SecondLevels = append(dt.SecondLevels, cate.Named)

		case thirdLevel[cate.Named] != nil:
			for key := range thirdLevel[cate.Named] {
				cate.Children = append(cate.Children, key)
			}
			// dt.ThirdLevels = append(dt.ThirdLevels, cate.Named)

		case fourthLevel[cate.Named] != nil:
			for key := range fourthLevel[cate.Named] {
				cate.Children = append(cate.Children, key)
			}
			// dt.FourthLevels = append(dt.FourthLevels, cate.Named)
		}
	}

	out := bytes.NewBufferString("")

	t := template.Must(template.New("t.go.tmpl").ParseFiles("t.go.tmpl"))
	if err = t.Execute(out, dt); err != nil {
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
