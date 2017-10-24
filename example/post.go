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
		Body: &request.Data{
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Status)
	fmt.Println(string(data))
}
