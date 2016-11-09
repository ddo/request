package request

import (
	"net/http"
	"testing"
)

var testing_cookieClient *Client

func TestCookieInit(t *testing.T) {
	testing_cookieClient = New()

	res, err := testing_cookieClient.Request(&Option{
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
	cookies, err := testing_cookieClient.GetCookies("http://httpbin.org")

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
	cookie2, err := testing_cookieClient.GetCookie("http://httpbin.org", "cookie2")

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

	testing_cookieClient.SetCookies("https://httpbin.org", cookies)

	res, err := testing_cookieClient.Request(&Option{
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
	cookies, _ := testing_cookieClient.GetCookies("http://httpbin.org")

	// override the old one
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "cookie2" {
			cookies[i].Value = "hai"
		}
	}

	testing_cookieClient.SetCookies("http://httpbin.org", cookies)

	res, err := testing_cookieClient.Request(&Option{
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
