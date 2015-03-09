package request

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

var c = New()

func TestNew(t *testing.T) {
	if c.httpClient == nil {
		t.Fail()
	}
}

func TestMakeUrl(t *testing.T) {
	u, err := makeUrl("https://httpbin.org/get?one=1", &Data{
		"two":   []string{"2", "hai"},
		"three": []string{"3", "ba", "trois"},
		"email": []string{"ddo@ddo.me"},
	})

	if err != nil {
		t.Fail()
	}

	if u.String() != "https://httpbin.org/get?email=ddo%40ddo.me&one=1&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Fail()
	}
}

func TestMakeUrlNil(t *testing.T) {
	u, err := makeUrl("https://httpbin.org/get?one=1", nil)

	if err != nil {
		t.Fail()
	}

	if u.String() != "https://httpbin.org/get?one=1" {
		t.Fail()
	}
}

func TestMakeBody(t *testing.T) {
	body := makeBody(&Option{
		Body: &Data{
			"one":   []string{"1", "mot"},
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})

	if body != "email=ddo%40ddo.me&one=1&one=mot&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Fail()
	}
}

func TestMakeBodyNil(t *testing.T) {
	body := makeBody(&Option{})

	if body != "" {
		t.Fail()
	}
}

func TestMakeBodyForm(t *testing.T) {
	body := makeBody(&Option{
		Form: &Data{
			"one":   []string{"1", "mot"},
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})

	if body != "email=ddo%40ddo.me&one=1&one=mot&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Fail()
	}
}

func TestMakeHeaderDefault(t *testing.T) {
	req, _ := http.NewRequest("POST", "https://httpbin.org", strings.NewReader(""))

	makeHeader(req, &Option{})

	if req.Header["User-Agent"][0] != "github.com/ddo/request" {
		t.Fail()
	}
}

func TestMakeHeaderForm(t *testing.T) {
	req, _ := http.NewRequest("POST", "https://httpbin.org", strings.NewReader(""))

	makeHeader(req, &Option{
		Form: &Data{},
	})

	if req.Header["User-Agent"][0] != "github.com/ddo/request" {
		t.Fail()
	}

	if req.Header["Content-Type"][0] != "application/x-www-form-urlencoded" {
		t.Fail()
	}
}

func TestMakeHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://httpbin.org", strings.NewReader(""))

	makeHeader(req, &Option{
		Header: &Header{
			"Custom":     "Custom header",
			"Custom2":    " ",
			"User-Agent": "",
		},
	})

	if req.Header["User-Agent"][0] != "" {
		t.Fail()
	}

	if req.Header["Custom"][0] != "Custom header" {
		t.Fail()
	}

	if req.Header["Custom2"][0] != " " {
		t.Fail()
	}
}

func TestRequest(t *testing.T) {
	client := New()

	body, res, err := client.Request(&Option{
		Url: "https://httpbin.org/ip",
	})

	if err != nil {
		t.Fail()
	}

	if res == nil {
		t.Fail()
	}

	if body == "" {
		t.Fail()
	}
}

func TestRequestRes(t *testing.T) {
	client := New()

	_, res, err := client.Request(&Option{
		Url: "https://httpbin.org/status/500",
	})

	if err != nil {
		t.Fail()
	}

	if res == nil {
		t.Fail()
	}

	if res.StatusCode != 500 {
		t.Fail()
	}
}

func TestRequestDefaultUserAgent(t *testing.T) {
	client := New()

	body, res, err := client.Request(&Option{
		Url: "https://httpbin.org/get",
	})

	if err != nil {
		t.Fail()
	}

	if res == nil {
		t.Fail()
	}

	if body == "" {
		t.Fail()
	}

	data := decodeHttpbinRes(body)

	userAgent := data.Headers["User-Agent"]

	if userAgent != "github.com/ddo/request" {
		t.Fail()
	}
}

func TestRequestHeader(t *testing.T) {
	client := New()

	body, res, err := client.Request(&Option{
		Url: "https://httpbin.org/get",
		Header: &Header{
			"Custom":     "Custom header",
			"User-Agent": "",
		},
	})

	if err != nil {
		t.Fail()
	}

	if res == nil {
		t.Fail()
	}

	if body == "" {
		t.Fail()
	}

	data := decodeHttpbinRes(body)

	if data.Headers["User-Agent"] != "" {
		t.Fail()
	}

	if data.Headers["Custom"] != "Custom header" {
		t.Fail()
	}
}

func TestRequestGET(t *testing.T) {
	client := New()

	body, res, err := client.Request(&Option{
		Url: "https://httpbin.org/get?one=1",
		Query: &Data{
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})

	if err != nil {
		t.Fail()
	}

	if res == nil {
		t.Fail()
	}

	if body == "" {
		t.Fail()
	}

	data := decodeHttpbinRes(body)

	email := data.Args["email"].(string)
	one := data.Args["one"].(string)
	two := data.Args["two"].([]interface{})
	three := data.Args["three"].([]interface{})

	if email != "ddo@ddo.me" {
		t.Fail()
	}

	if one != "1" {
		t.Fail()
	}

	if two[0].(string) != "2" {
		t.Fail()
	}

	if two[1].(string) != "hai" {
		t.Fail()
	}

	if three[0].(string) != "3" {
		t.Fail()
	}

	if three[1].(string) != "ba" {
		t.Fail()
	}

	if three[2].(string) != "trois" {
		t.Fail()
	}
}

func TestRequestPOSTStr(t *testing.T) {
	client := New()

	body, res, err := client.Request(&Option{
		Url:    "https://httpbin.org/post",
		Method: "POST",
		Query: &Data{
			"one": []string{"1"},
		},
		BodyStr: "email=ddo%40ddo.me&three=3&three=ba&three=trois&two=2&two=hai",
	})

	if err != nil {
		t.Fail()
	}

	if res == nil {
		t.Fail()
	}

	if body == "" {
		t.Fail()
	}

	data := decodeHttpbinRes(body)

	if data.Args["one"].(string) != "1" {
		t.Fail()
	}

	if data.Data != "email=ddo%40ddo.me&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Fail()
	}
}

func TestRequestPOST(t *testing.T) {
	client := New()

	body, res, err := client.Request(&Option{
		Url:    "https://httpbin.org/post",
		Method: "POST",
		Query: &Data{
			"one": []string{"1"},
		},
		Body: &Data{
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})

	if err != nil {
		t.Fail()
	}

	if res == nil {
		t.Fail()
	}

	if body == "" {
		t.Fail()
	}

	data := decodeHttpbinRes(body)

	if data.Args["one"].(string) != "1" {
		t.Fail()
	}

	if data.Data != "email=ddo%40ddo.me&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Fail()
	}
}

func TestRequestPOSTForm(t *testing.T) {
	client := New()

	body, res, err := client.Request(&Option{
		Url:    "https://httpbin.org/post",
		Method: "post",
		Query: &Data{
			"one": []string{"1"},
		},
		Form: &Data{
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})

	if err != nil {
		t.Fail()
	}

	if res == nil {
		t.Fail()
	}

	if body == "" {
		t.Fail()
	}

	data := decodeHttpbinRes(body)

	if data.Args["one"].(string) != "1" {
		t.Fail()
	}

	email := data.Form["email"].(string)
	two := data.Form["two"].([]interface{})
	three := data.Form["three"].([]interface{})

	if email != "ddo@ddo.me" {
		t.Fail()
	}

	if two[0].(string) != "2" {
		t.Fail()
	}

	if two[1].(string) != "hai" {
		t.Fail()
	}

	if three[0].(string) != "3" {
		t.Fail()
	}

	if three[1].(string) != "ba" {
		t.Fail()
	}

	if three[2].(string) != "trois" {
		t.Fail()
	}
}

func TestRequestFail(t *testing.T) {
	client := New()

	body, res, err := client.Request(&Option{
		Url: "http://1.com",
	})

	if err == nil {
		t.Fail()
	}

	if res != nil {
		t.Fail()
	}

	if body != "" {
		t.Fail()
	}
}

////// helper

type httpbinRes struct {
	Args    map[string]interface{} `json:args`
	Headers map[string]string      `json:headers`
	Data    string                 `json:data`
	Form    map[string]interface{} `json:form`
}

func decodeHttpbinRes(body string) *httpbinRes {
	var data httpbinRes
	json.Unmarshal([]byte(body), &data)
	return &data
}
