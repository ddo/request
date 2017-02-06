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
	}

	if res == nil {
		t.Error()
	}

	data := decodeHttpbinRes(res)

	if data.Cookies["cookie1"] != "one" {
		t.Error()
	}

	if data.Cookies["cookie2"] != "2" {
		t.Error()
	}
}

func TestGetCookies(t *testing.T) {
	cookies, err := testCookieClient.GetCookies("http://httpbin.org")

	if err != nil {
		t.Error()
	}

	if len(cookies) != 2 {
		t.Error()
	}

	if cookies[0].Name != "cookie1" {
		t.Error()
	}

	if cookies[0].Value != "one" {
		t.Error()
	}

	if cookies[1].Name != "cookie2" {
		t.Error()
	}

	if cookies[1].Value != "2" {
		t.Error()
	}
}

func TestGetCookie(t *testing.T) {
	cookie2, err := testCookieClient.GetCookie("http://httpbin.org", "cookie2")

	if err != nil {
		t.Error()
	}

	if cookie2 != "2" {
		t.Error()
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

		// Name  string
		// Value string

		// Path       string    // optional
		// Domain     string    // optional
		// Expires    time.Time // optional
		// RawExpires string    // for reading cookies only

		// // MaxAge=0 means no 'Max-Age' attribute specified.
		// // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
		// // MaxAge>0 means Max-Age attribute present and given in seconds
		// MaxAge   int
		// Secure   bool
		// HttpOnly bool
		// Raw      string
		// Unparsed []string
	}

	cookies = append(cookies, newCookie)

	testCookieClient.SetCookies("https://httpbin.org", cookies)

	res, err := testCookieClient.Request(&Option{
		Url: "http://httpbin.org/cookies",
	})

	if err != nil {
		t.Error()
	}

	if res == nil {
		t.Error()
	}

	data := decodeHttpbinRes(res)

	if data.Cookies["cookie3"] != "ba" {
		t.Error()
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
	}

	if res == nil {
		t.Error()
	}

	data := decodeHttpbinRes(res)

	if data.Cookies["cookie2"] != "hai" {
		t.Error()
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
	}

	if data.Cookies["cookie2"] != "2" {
		t.Error()
	}
}
