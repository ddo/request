# request [![Build Status][travis-img]][travis-url] [![Doc][godoc-img]][godoc-url]
> Simplified HTTP request client in go

[travis-img]: https://img.shields.io/travis/ddo/request.svg?style=flat-square
[travis-url]: https://travis-ci.org/ddo/request

[godoc-img]: https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat-square
[godoc-url]: https://godoc.org/github.com/ddo/request

##Options

* Url     ``string`` required
* Method  ``string`` default: "GET", anything "POST", "PUT", "DELETE" or "PATCH"
* BodyStr ``string``
* Body    ``*Data``
* Form    ``*Data``       set Content-Type header as "application/x-www-form-urlencoded"
* Json    ``interface{}`` set Content-Type header as "application/json"
* Query   ``*Data``
* Header  ``*Header``

###GET

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

###POST

```go
body, res, err := client.Request(&request.Option{
    Url:    "https://httpbin.org/post",
    Method: "POST",
    Body: &request.Data{
        "two":   []string{"2", "hai"},
        "three": []string{"3", "ba", "trois"},
        "email": []string{"ddo@ddo.me"},
    },
})
```

###POST form

```go
body, res, err := client.Request(&request.Option{
    Url:    "https://httpbin.org/post",
    Method: "POST",
    Form: &request.Data{
        "two":   []string{"2", "hai"},
        "three": []string{"3", "ba", "trois"},
        "email": []string{"ddo@ddo.me"},
    },
})
```

###Json

```go
body, res, err := client.Request(&request.Option{
    Url:    "https://httpbin.org/post",
    Method: "POST",
    Json: map[string]interface{}{
        "int":    1,
        "string": "two",
        "array":  []string{"3", "ba", "trois"},
        "object": map[string]interface{}{
            "int": 4,
        },
    },
})
```

##Logger

to enable log set environment variable as

```go
DEBUG=request
```

or

```shell
DEBUG=request go run file.go
```

##Test

```shell
go test -v
```

##TODO

* default settings
* hooks
* file