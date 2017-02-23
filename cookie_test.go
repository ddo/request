package request

import (
	"net/http"
	"testing"
)

var testCookieClient *Client

func TestCookieInit(t *testing.T) {
	testCookieClient = New()

	res, err := testCookieClient.Request(&Option{
		Url: "http://httpbin.org/cookies/set?cookie1=one&cookie2=2",
	})
	if err != nil {
		t.Error()
		return
	}

	if res == nil {
		t.Error()
		return
	}

	data := decodeHttpbinRes(res)

	if data.Cookies["cookie1"] != "one" {
		t.Error()
		return
	}

	if data.Cookies["cookie2"] != "2" {
		t.Error()
		return
	}
}

func TestGetCookies(t *testing.T) {
	cookies, err := testCookieClient.GetCookies("http://httpbin.org")
	if err != nil {
		t.Error()
		return
	}

	if len(cookies) != 2 {
		t.Error()
		return
	}

	if cookies[0].Name != "cookie1" {
		t.Error()
		return
	}

	if cookies[0].Value != "one" {
		t.Error()
		return
	}

	if cookies[1].Name != "cookie2" {
		t.Error()
		return
	}

	if cookies[1].Value != "2" {
		t.Error()
		return
	}
}

func TestGetCookie(t *testing.T) {
	cookie2, err := testCookieClient.GetCookie("http://httpbin.org", "cookie2")
	if err != nil {
		t.Error()
		return
	}

	if cookie2 != "2" {
		t.Error()
		return
	}
}

func TestSetCookies(t *testing.T) {
	// empty cookie slice
	cookies := []*http.Cookie{}

	// new cookie
	newCookie := &http.Cookie{
		Name:   "cookie3",
		Value:  "ba",
		MaxAge: 0,
	}

	cookies = append(cookies, newCookie)

	testCookieClient.SetCookies("https://httpbin.org", cookies)

	res, err := testCookieClient.Request(&Option{
		Url: "http://httpbin.org/cookies",
	})
	if err != nil {
		t.Error()
		return
	}

	if res == nil {
		t.Error()
		return
	}

	data := decodeHttpbinRes(res)

	if data.Cookies["cookie3"] != "ba" {
		t.Error()
		return
	}
}

func TestSetCookiesModify(t *testing.T) {
	cookies, _ := testCookieClient.GetCookies("http://httpbin.org")

	// override the old one
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "cookie2" {
			cookies[i].Value = "hai"
		}
	}

	testCookieClient.SetCookies("http://httpbin.org", cookies)

	res, err := testCookieClient.Request(&Option{
		Url: "https://httpbin.org/cookies",
	})
	if err != nil {
		t.Error()
		return
	}

	if res == nil {
		t.Error()
		return
	}

	data := decodeHttpbinRes(res)

	if data.Cookies["cookie2"] != "hai" {
		t.Error()
		return
	}
}

func TestImportCookie(t *testing.T) {
	client := New()
	err := client.ImportCookie("http://httpbin.org", `
	[{
		"domain": ".twitter.com",
		"expires": "Wed, 06 Feb 2019 09:29:47 GMT",
		"expiry": 1549445387,
		"httponly": false,
		"name": "guest_id",
		"path": "/",
		"secure": false,
		"value": "v1%3A1486379578632445354"
	}, {
		"domain": "www.abc.com",
		"httponly": false,
		"name": "QSI_HistorySession",
		"path": "/",
		"secure": false,
		"value": "http%3A%2F%2Fwww.abc.com%2Fus~1486373387091"
	}, {
		"domain": ".httpbin.org",
		"httponly": false,
		"name": "cookie1",
		"path": "/",
		"secure": false,
		"value": "1"
	}, {
		"domain": "httpbin.org",
		"httponly": false,
		"name": "cookie2",
		"path": "/",
		"secure": false,
		"value": "2"
	}]`)
	if err != nil {
		t.Error()
		return
	}

	// verify
	res, err := client.Request(&Option{
		Url: "https://httpbin.org/cookies",
	})
	if err != nil {
		t.Error()
		return
	}
	if res == nil {
		t.Error()
		return
	}

	data := decodeHttpbinRes(res)

	if data.Cookies["cookie1"] != "1" {
		t.Error()
		return
	}

	if data.Cookies["cookie2"] != "2" {
		t.Error()
		return
	}
}

func TestExportCookie(t *testing.T) {
	// create client with cookies
	client := New()

	res, err := client.Request(&Option{
		Url: "http://httpbin.org/cookies/set?cookie1=1&cookie2=2",
	})
	if err != nil {
		t.Error()
		return
	}
	defer res.Body.Close()
	// create client with cookies - end

	jsonStr, err := client.ExportCookie("http://httpbin.org")
	if err != nil {
		t.Error()
		return
	}

	sample := `[{"name":"cookie1","value":"1","path":"","domain":"","secure":false,"httponly":false},{"name":"cookie2","value":"2","path":"","domain":"","secure":false,"httponly":false}]`
	sampleCustomCookieJar := `[{"name":"cookie1","value":"1","path":"/","domain":"httpbin.org","secure":false,"httponly":false},{"name":"cookie2","value":"2","path":"/","domain":"httpbin.org","secure":false,"httponly":false}]`

	// NOTE: local test on a customize cookie jar
	if jsonStr != sample && jsonStr != sampleCustomCookieJar {
		t.Error(jsonStr)
		return
	}
}
