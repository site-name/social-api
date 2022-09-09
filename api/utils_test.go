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
	b := person{"minh", 24}

	c := context.WithValue(context.Background(), WebCtx, &b)

	p, err := GetContextValue[person](c, WebCtx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p.age, p.name)
}
