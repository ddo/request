package request

// TODO: add tests
// GetCookies	71.4%
// SetCookies	75.0%
// GetCookie	69.2%
// New			100.0%
// NewNoCookie	0.0%
// SetTimeout	0.0%
// SetProxy		0.0%
// Request		81.8%
// makeURL		83.3%
// makeBody		89.5%
// makeHeader	100.0%

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	defaultHeader = " "
	streamLength  = 50
)

var testClient *Client

func init() {
	testClient = New()
}

func TestNew(t *testing.T) {
	if testClient.httpClient == nil {
		t.Error()
		return
	}
}

func TestMakeURL(t *testing.T) {
	u, err := makeURL("https://httpbin.org/get?one=1", &Data{
		"two":   []string{"2", "hai"},
		"three": []string{"3", "ba", "trois"},
		"email": []string{"ddo@ddo.me"},
	})
	if err != nil {
		t.Error()
		return
	}

	if u.String() != "https://httpbin.org/get?email=ddo%40ddo.me&one=1&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Error()
		return
	}
}

func TestMakeURLNil(t *testing.T) {
	u, err := makeURL("https://httpbin.org/get?one=1", nil)
	if err != nil {
		t.Error()
		return
	}

	if u.String() != "https://httpbin.org/get?one=1" {
		t.Error()
		return
	}
}

func TestMakeBody(t *testing.T) {
	body, err := makeBody(&Option{
		Body: &Data{
			"one":   []string{"1", "mot"},
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})
	if err != nil {
		t.Error()
		return
	}

	if body != "email=ddo%40ddo.me&one=1&one=mot&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Error()
		return
	}
}

func TestMakeBodyNil(t *testing.T) {
	body, err := makeBody(&Option{})
	if err != nil {
		t.Error()
		return
	}

	if body != "" {
		t.Error()
		return
	}
}

func TestMakeBodyForm(t *testing.T) {
	body, err := makeBody(&Option{
		Form: &Data{
			"one":   []string{"1", "mot"},
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})
	if err != nil {
		t.Error()
		return
	}

	if body != "email=ddo%40ddo.me&one=1&one=mot&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Error()
		return
	}
}

func TestMakeBodyJson(t *testing.T) {
	body, err := makeBody(&Option{
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
		t.Error()
		return
	}

	if body != `{"array":["3","ba","trois"],"int":1,"object":{"int":4},"string":"two"}` {
		t.Error()
		return
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
		t.Error()
		return
	}

	if req.Header["Custom"][0] != "Custom header" {
		t.Error()
		return
	}

	if req.Header["Custom2"][0] != " " {
		t.Error()
		return
	}
}

func TestMakeHeaderDefault(t *testing.T) {
	req, _ := http.NewRequest("POST", "https://httpbin.org", strings.NewReader(""))

	makeHeader(req, &Option{})

	if req.Header["User-Agent"][0] != defaultHeader {
		t.Error()
		return
	}
}

func TestMakeHeaderForm(t *testing.T) {
	req, _ := http.NewRequest("POST", "https://httpbin.org", strings.NewReader(""))

	makeHeader(req, &Option{
		Form: &Data{},
	})

	if req.Header["User-Agent"][0] != defaultHeader {
		t.Error()
		return
	}

	if req.Header["Content-Type"][0] != "application/x-www-form-urlencoded" {
		t.Error()
		return
	}
}

func TestMakeHeaderJson(t *testing.T) {
	req, _ := http.NewRequest("POST", "https://httpbin.org", strings.NewReader(""))

	makeHeader(req, &Option{
		JSON: &Data{},
	})

	if req.Header["User-Agent"][0] != defaultHeader {
		t.Error()
		return
	}

	if req.Header["Content-Type"][0] != "application/json" {
		t.Error()
		return
	}
}

func TestRequest(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL: "https://httpbin.org/ip",
	})
	if err != nil {
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}
}

func TestRequestRes(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL: "https://httpbin.org/status/500",
	})
	if err != nil {
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	if res.StatusCode != 500 {
		t.Error()
		return
	}
}

func TestRequestDefaultUserAgent(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL: "https://httpbin.org/get",
	})
	if err != nil {
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	testData := decodeHttpbinRes(data)

	userAgent := testData.Headers["User-Agent"]

	if userAgent != "" {
		t.Error()
		return
	}
}

func TestRequestHeader(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL: "https://httpbin.org/get",
		Header: &Header{
			"Custom":     "Custom header",
			"User-Agent": "",
		},
	})
	if err != nil {
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	testData := decodeHttpbinRes(data)

	if testData.Headers["User-Agent"] != "" {
		t.Error()
		return
	}

	if testData.Headers["Custom"] != "Custom header" {
		t.Error()
		return
	}
}

func TestRequestGET(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL: "https://httpbin.org/get?one=1",
		Query: &Data{
			"two":   []string{"2", "hai"},
			"three": []string{"3", "ba", "trois"},
			"email": []string{"ddo@ddo.me"},
		},
	})
	if err != nil {
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	testData := decodeHttpbinRes(data)

	email := testData.Args["email"].(string)
	one := testData.Args["one"].(string)
	two := testData.Args["two"].([]interface{})
	three := testData.Args["three"].([]interface{})

	if email != "ddo@ddo.me" {
		t.Error()
		return
	}

	if one != "1" {
		t.Error()
		return
	}

	if two[0].(string) != "2" {
		t.Error()
		return
	}

	if two[1].(string) != "hai" {
		t.Error()
		return
	}

	if three[0].(string) != "3" {
		t.Error()
		return
	}

	if three[1].(string) != "ba" {
		t.Error()
		return
	}

	if three[2].(string) != "trois" {
		t.Error()
		return
	}
}

func TestRequestPOSTStr(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL:    "https://httpbin.org/post",
		Method: "POST",
		Query: &Data{
			"one": []string{"1"},
		},
		BodyStr: "email=ddo%40ddo.me&three=3&three=ba&three=trois&two=2&two=hai",
	})
	if err != nil {
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	testData := decodeHttpbinRes(data)

	if testData.Args["one"].(string) != "1" {
		t.Error()
		return
	}

	if testData.Data != "email=ddo%40ddo.me&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Error()
		return
	}
}

func TestRequestPOST(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL:    "https://httpbin.org/post",
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
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	testData := decodeHttpbinRes(data)

	if testData.Args["one"].(string) != "1" {
		t.Error()
		return
	}

	if testData.Data != "email=ddo%40ddo.me&three=3&three=ba&three=trois&two=2&two=hai" {
		t.Error()
		return
	}
}

func TestRequestPOSTJson(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL:    "https://httpbin.org/post",
		Method: "POST",
		Query: &Data{
			"one": []string{"1"},
		},
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
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	testData := decodeHttpbinRes(data)

	if testData.Args["one"].(string) != "1" {
		t.Error()
		return
	}

	if testData.Data != `{"array":["3","ba","trois"],"int":1,"object":{"int":4},"string":"two"}` {
		t.Error()
		return
	}

	if testData.JSON.Int != 1 {
		t.Error()
		return
	}

	if testData.JSON.String != "two" {
		t.Error()
		return
	}

	if testData.JSON.Array[0] != "3" {
		t.Error()
		return
	}

	if testData.JSON.Array[1] != "ba" {
		t.Error()
		return
	}

	if testData.JSON.Array[2] != "trois" {
		t.Error()
		return
	}

	if testData.JSON.Object["int"] != 4 {
		t.Error()
		return
	}
}

func TestRequestPOSTForm(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL:    "https://httpbin.org/post",
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
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	testData := decodeHttpbinRes(data)

	if testData.Args["one"].(string) != "1" {
		t.Error()
		return
	}

	email := testData.Form["email"].(string)
	two := testData.Form["two"].([]interface{})
	three := testData.Form["three"].([]interface{})

	if email != "ddo@ddo.me" {
		t.Error()
		return
	}

	if two[0].(string) != "2" {
		t.Error()
		return
	}

	if two[1].(string) != "hai" {
		t.Error()
		return
	}

	if three[0].(string) != "3" {
		t.Error()
		return
	}

	if three[1].(string) != "ba" {
		t.Error()
		return
	}

	if three[2].(string) != "trois" {
		t.Error()
		return
	}
}

func TestRequestFail(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL: "http://1.com",
	})
	if err == nil {
		t.Error()
		return
	}
	if res != nil {
		t.Error()
		return
	}
	if data != nil {
		t.Error()
		return
	}
}

func TestRequestStream(t *testing.T) {
	client := New()

	data, res, err := client.Request(&Option{
		URL: fmt.Sprintf("https://httpbin.org/stream/%v", streamLength),
	})
	if err != nil {
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}
	if data == nil {
		t.Error()
		return
	}

	// process stream
	counter := 0

	var testData httpbinRes
	decoder := json.NewDecoder(bytes.NewReader(data))

	for {
		err := decoder.Decode(&testData)

		if err == io.EOF {
			break
		}

		// ignore the error
		if err != nil {
			panic(err)
		}

		counter++
	}

	if counter != streamLength {
		t.Error()
		return
	}
}

////// helper

type httpbinResJSON struct {
	Int    int            `json:"int"`
	String string         `json:"string"`
	Array  []string       `json:"array"`
	Object map[string]int `json:"object"`
}

type httpbinRes struct {
	Args    map[string]interface{} `json:"args"`
	Headers map[string]string      `json:"headers"`
	Data    string                 `json:"data"`
	Form    map[string]interface{} `json:"form"`
	JSON    httpbinResJSON         `json:"json"`
	Cookies map[string]string      `json:"cookies"`
}

func decodeHttpbinRes(data []byte) *httpbinRes {
	// debug(string(data))

	var res httpbinRes
	json.Unmarshal(data, &res)
	return &res
}
