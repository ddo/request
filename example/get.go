package main

import (
	"fmt"

	"github.com/ddo/request"
)

func main() {
	client := request.New()

	body, res, err := client.Request(&request.Option{
		Url: "https://httpbin.org/get?one=1",
		Query: &request.Query{
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	fmt.Println(body)
}
