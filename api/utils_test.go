package api

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/sitename/sitename/model"
)

type person struct {
	name string
	age  int
}

func TestGetContextValue(t *testing.T) {
	b := person{"minh", 24}

	c := context.WithValue(context.Background(), WebCtx, &b)

	p, err := GetContextValue[person](c, WebCtx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p.age, p.name)
}

func TestGraphqlPaginator(t *testing.T) {
	type person struct {
		name string
		age  int
	}

	var persons = []person{
		{"minh", 22},
		{"fggghhngnbn", 34},
		{"cvfdgtg", 3},
		{"cvghnhgnh", 22},
		{"saertreggb", 44},
		{"nhjio7e6fsfa", 55},
		{"bgbaeraghht", 16},
	}

	p := graphqlPaginator[person, string]{
		data:    persons,
		keyFunc: func(p person) string { return p.name },
		first:   model.NewPrimitive[int32](1),
		after:   model.NewPrimitive("minh"),
	}

	data, hasPreviousPage, hasNextPage, err := p.parse("Something")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(p.data)

	fmt.Println(data, hasPreviousPage, hasNextPage)
}
