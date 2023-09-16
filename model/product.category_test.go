package model

// func TestClassifyCategories(t *testing.T) {
// 	in := Categories{
// 		&Category{Id: "1", ParentID: nil},
// 		&Category{Id: "2", ParentID: GetPointerOfValue("1")},
// 		&Category{Id: "3", ParentID: GetPointerOfValue("2")},
// 		&Category{Id: "4", ParentID: GetPointerOfValue("3")},
// 		&Category{Id: "5", ParentID: GetPointerOfValue("4")},
// 		&Category{Id: "6", ParentID: nil},
// 		&Category{Id: "7", ParentID: GetPointerOfValue("1")},
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
