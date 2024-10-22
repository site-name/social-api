package api

import (
	"context"
	"fmt"
	"testing"
)

type person struct {
	name string
	age  int
}

func TestGetContextValue(t *testing.T) {
	minh := person{"minh", 24}

	c := context.WithValue(context.Background(), WebCtx, minh)

	p := GetContextValue[person](c, WebCtx)

	fmt.Println(p.age, p.name)
}

// func TestGraphqlPaginator(t *testing.T) {
// 	type person struct {
// 		name string
// 		age  int
// 	}

// 	var persons = []person{
// 		{"minh", 22},
// 		{"fggghhngnbn", 34},
// 		{"cvfdgtg", 3},
// 		{"cvghnhgnh", 22},
// 		{"saertreggb", 44},
// 		{"nhjio7e6fsfa", 55},
// 		{"bgbaeraghht", 16},
// 	}

// 	p := graphqlPaginator[person, string]{
// 		data:    persons,
// 		keyFunc: func(p person) string { return p.name },
// 		last:    model_helper.GetPointerOfValue[int32](1),
// 		before:  model_helper.GetPointerOfValue("cvfdgtg"),
// 	}

// 	data, hasPreviousPage, hasNextPage, err := p.parse("Something")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	fmt.Println(p.data)

// 	fmt.Println(data, hasPreviousPage, hasNextPage)
// }
