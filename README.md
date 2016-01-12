# request [![Build Status][semaphoreci-img]][semaphoreci-url] [![Doc][godoc-img]][godoc-url]
> Simplified HTTP request client in go

[godoc-img]: https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat-square
[godoc-url]: https://godoc.org/gopkg.in/ddo/request.v1

[semaphoreci-img]: https://semaphoreci.com/api/v1/projects/fe48ba6a-f987-4018-b778-34c0fef12c87/620801/badge.svg
[semaphoreci-url]: https://semaphoreci.com/ddo/request

## installation
```sh
go get gopkg.in/ddo/request.v1
```

## option

* Url     ``string`` required
* Method  ``string`` default: "GET", anything "POST", "PUT", "DELETE" or "PATCH"
* BodyStr ``string``
* Body    ``*Data``
* Form    ``*Data``       set Content-Type header as "application/x-www-form-urlencoded"
* Json    ``interface{}`` set Content-Type header as "application/json"
* Query   ``*Data``
* Header  ``*Header``

### GET

```go
client := request.New()

res, err := client.Request(&request.Option{
    Url: "https://httpbin.org/get",
})

if err != nil {
    panic(err)
}

defer res.Body.Close()
```

### POST

```go
res, err := client.Request(&request.Option{
    Url:    "https://httpbin.org/post",
    Method: "POST",
    Body: &request.Data{
        "two":   []string{"2", "hai"},
        "three": []string{"3", "ba", "trois"},
        "email": []string{"ddo@ddo.me"},
    },
})
```

### POST form

```go
res, err := client.Request(&request.Option{
    Url:    "https://httpbin.org/post",
    Method: "POST",
    Form: &request.Data{
        "two":   []string{"2", "hai"},
        "three": []string{"3", "ba", "trois"},
        "email": []string{"ddo@ddo.me"},
    },
})
```

### Json

```go
res, err := client.Request(&request.Option{
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

## logger

to enable log set environment variable as

```go
DLOG=*
```

or

```shell
DEBUG=* go run file.go
```

## test

```shell
go test -v
```

## TODO

* default settings
* hooks
* file