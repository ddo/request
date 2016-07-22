package request

import (
	"net/http"
	"net/url"
)

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
