package request

import (
	"net/http"
	"net/url"
)

// SetProxy sets client proxy
func (c *Client) SetProxy(proxyURLStr string) (err error) {
	debug(proxyURLStr)

	proxyURL, err := url.Parse(proxyURLStr)
	if err != nil {
		debug("ERR(url.Parse)", err)
		return
	}

	c.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	return
}
