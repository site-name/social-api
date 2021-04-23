package main

import (
	"fmt"
	"net/url"
)

func main() {
	vl := url.Values{}

	vl["name"] = []string{"Konomi Nagisa"}
	vl["age"] = []string{"12"}

	res := vl.Encode()

	fmt.Println(res)
}
