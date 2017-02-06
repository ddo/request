package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GetCookies gets []*http.Cookie by domain
func (c *Client) GetCookies(domain string) (cookies []*http.Cookie, err error) {
	debug(domain)

	u, err := url.Parse(domain)
	if err != nil {
		debug("ERR(parse)", err)
		return
	}

	cookies = c.httpClient.Jar.Cookies(u)
	return
}

// SetCookies sets cookies by domain
func (c *Client) SetCookies(domain string, cookies []*http.Cookie) (err error) {
	debug(domain)
	debug(cookies)

	u, err := url.Parse(domain)
	if err != nil {
		debug("ERR(parse)", err)
		return
	}

	c.httpClient.Jar.SetCookies(u, cookies)
	return
}

// GetCookie gets cookie value by domain and cookie name
func (c *Client) GetCookie(domain, name string) (value string, err error) {
	debug(domain, name)

	u, err := url.Parse(domain)
	if err != nil {
		debug("ERR(parse)", err)
		return
	}

	cookies := c.httpClient.Jar.Cookies(u)

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == name {
			value = cookies[i].Value

			debug("DONE", value)
			return
		}
	}

	debug("EMPTY")
	return
}

// Cookie is a duplicate of http.Cookie but with json parser
type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path       string `json:"path"`
	Domain     string `json:"domain"`
	Expires    time.Time
	RawExpires string `json:"expires"` // for reading cookies only

	MaxAge   int  `json:"expiry"`
	Secure   bool `json:"secure"`
	HttpOnly bool `json:"httponly"`
}

// debug purpose
func (c Cookie) String() string {
	return fmt.Sprintf("\nName\t\t:%s\nValue\t\t:%s\nPath\t\t:%s\nDomain\t\t:%s\nExpires\t\t:%v\nRawExpires\t:%s\nMaxAge\t\t:%v\nSecure\t\t:%v\nHttpOnly\t:%v\n-------------\n", c.Name, c.Value, c.Path, c.Domain, c.Expires, c.RawExpires, c.MaxAge, c.Secure, c.HttpOnly)
}

// ImportCookie imports cookie from json
func (c *Client) ImportCookie(domain, jsonStr string) (err error) {
	debug("domain:", domain)

	var cookies []Cookie

	err = json.Unmarshal([]byte(jsonStr), &cookies)
	if err != nil {
		debug("ERR(json.Unmarshal)", err)
		return
	}

	// debug(cookies)

	httpCookies := tohttpCookie(cookies)

	err = c.SetCookies(domain, httpCookies)
	if err != nil {
		return
	}
	return
}

func tohttpCookie(cookies []Cookie) (httpCookies []*http.Cookie) {
	debug()

	var expires time.Time
	var err error

	for i := 0; i < len(cookies); i++ {
		// .Expires
		expires, err = time.Parse(time.RFC1123, cookies[i].RawExpires)
		if err == nil {
			cookies[i].Expires = expires
		}

		// new httpCookie
		httpCookies = append(httpCookies, &http.Cookie{
			Name:       cookies[i].Name,
			Value:      cookies[i].Value,
			Path:       cookies[i].Path,
			Domain:     cookies[i].Domain,
			Expires:    cookies[i].Expires,
			RawExpires: cookies[i].RawExpires,
			MaxAge:     cookies[i].MaxAge,
			Secure:     cookies[i].Secure,
			HttpOnly:   cookies[i].HttpOnly,
		})
	}

	return
}

// TODO: implement ImportCookie
