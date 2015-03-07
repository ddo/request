package main

import (
	"fmt"

	"github.com/ddo/request"
)

func main() {
	client := request.New()

	body, res, err := client.Request(&request.Option{
		Url:    "https://httpbin.org/post",
		Method: "POST",
		Form: &request.Form{
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
		Header: &request.Header{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	fmt.Println(body)
}