package main

import (
	"fmt"

	"github.com/ddo/request"
)

func main() {
	client := request.New()

	data, res, err := client.Request(&request.Option{
		URL:    "https://httpbin.org/post",
		Method: "POST",
		JSON: map[string]interface{}{
			"int":    1,
			"string": "two",
			"array":  []string{"3", "ba", "trois"},
			"object": map[string]interface{}{
				"int": 4,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Status)
	fmt.Println(string(data))
}
