package request

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	. "github.com/tj/go-debug"
)

var debug = Debug("request")

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	var cookie, _ = cookiejar.New(nil)

	client := &http.Client{
		Timeout: time.Second * 30,
		Jar:     cookie,
	}

	debug("#New")

	return &Client{client}
}

type Data map[string][]string
type Header map[string]string

type Option struct {
	Url     string
	Method  string
	BodyStr string
	Body    *Data
	Form    *Data
	Query   *Data
	Header  *Header
}

func (c *Client) Request(opt *Option) (body string, res *http.Response, err error) {
	debug("#Request")

	//set GET as default method
	if opt.Method == "" {
		opt.Method = "GET"
	}

	opt.Method = strings.ToUpper(opt.Method)

	//url
	reqUrl, err := makeUrl(opt.Url, opt.Query)

	if err != nil {
		return
	}

	//body
	reqBody := makeBody(opt)

	req, err := http.NewRequest(opt.Method, reqUrl.String(), strings.NewReader(reqBody))

	if err != nil {
		debug("#Request ERR(req) %v", err)
		return
	}

	//header
	makeHeader(req, opt.Header)

	res, err = c.httpClient.Do(req)

	if err != nil {
		debug("#Request ERR(http) %v", err)
		return
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		debug("#Request ERR(ioutil) %v", err)
		return
	}

	body = string(resBody)

	debug("#Request %v", res.Status)
	return
}

func makeUrl(urlStr string, query *Data) (u *url.URL, err error) {
	// debug("#makeUrl")

	u, err = url.Parse(urlStr)

	if err != nil {
		debug("#makeUrl ERR: %v", err)
		return
	}

	if query == nil {
		return
	}

	qs := u.Query()

	for key, slice := range *query {
		for _, value := range slice {
			qs.Add(key, value)
		}
	}

	u.RawQuery = qs.Encode()
	return
}

func makeBody(opt *Option) (body string) {
	var data *Data

	switch {
	case opt.BodyStr != "":
		body = opt.BodyStr
		return

	case opt.Form != nil:
		data = opt.Form

	case opt.Body != nil:
		data = opt.Body

	default:
		return
	}

	values := url.Values{}

	for key, slice := range *data {
		for _, value := range slice {
			values.Add(key, value)
		}
	}

	body = values.Encode()
	return
}

func makeHeader(req *http.Request, header *Header) {
	// debug("#makeHeader")

	//default User-Agent
	req.Header.Set("User-Agent", "github.com/ddo/request")

	if header == nil {
		return
	}

	for key, value := range *header {
		req.Header.Set(key, value)
	}
}
