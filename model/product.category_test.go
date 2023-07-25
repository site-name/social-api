package model

// func TestClassifyCategories(t *testing.T) {
// 	in := Categories{
// 		&Category{Id: "1", ParentID: nil},
// 		&Category{Id: "2", ParentID: NewPrimitive("1")},
// 		&Category{Id: "3", ParentID: NewPrimitive("2")},
// 		&Category{Id: "4", ParentID: NewPrimitive("3")},
// 		&Category{Id: "5", ParentID: NewPrimitive("4")},
// 		&Category{Id: "6", ParentID: nil},
// 		&Category{Id: "7", ParentID: NewPrimitive("1")},
// 	}

// 	t.Run("ello", func(t *testing.T) {
// 		out := ClassifyCategories(in)
// 		data, err := json.MarshalIndent(out, " ", "  ")
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		fmt.Println(string(data))
// 	})
// }
