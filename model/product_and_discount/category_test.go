package product_and_discount

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sitename/sitename/model"
)

func TestClassifyCategories(t *testing.T) {
	in := Categories{
		&Category{Id: "1", ParentID: nil},
		&Category{Id: "2", ParentID: model.NewString("1")},
		&Category{Id: "3", ParentID: model.NewString("2")},
		&Category{Id: "4", ParentID: model.NewString("3")},
		&Category{Id: "5", ParentID: model.NewString("4")},
		&Category{Id: "6", ParentID: nil},
		&Category{Id: "7", ParentID: model.NewString("1")},
	}

	t.Run("ello", func(t *testing.T) {
		out := ClassifyCategories(in)
		data, err := json.MarshalIndent(out, " ", "  ")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(data))
	})
}
