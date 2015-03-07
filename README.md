# request [![Build Status][travis-img]][travis-url]
> Simplified HTTP request client in go

[![Doc][godoc-img]][godoc-url]

[travis-img]: https://img.shields.io/travis/ddo/request.svg?style=flat-square
[travis-url]: https://travis-ci.org/ddo/request

[godoc-img]: https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat-square
[godoc-url]: https://godoc.org/github.com/ddo/request

##GET

```go
client := request.New()

body, res, err := client.Request(&request.Option{
    Url: "https://httpbin.org/get",
})

if err != nil {
    panic(err)
}

fmt.Println(res)
fmt.Println(body)
```

##POST

```go
client := request.New()

body, res, err := client.Request(&request.Option{
    Url:    "https://httpbin.org/post",
    Method: "POST",
    Form: &request.Form{
        "two":   []string{"2", "hai"},
        "three": []string{"3", "ba", "trois"},
        "email": []string{"ddo@ddo.me"},
    },
})
```

##POST form

```go
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
```