package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ddo/go-dlog.v1"
)

const (
	// DefaultTimeout is request timeout in second
	DefaultTimeout = 180
)

var debug = dlog.New("request", nil)

// Client is an http client that hold init settings and cookies
type Client struct {
	httpClient *http.Client
}

// New return a new Client
func New() *Client {
	var cookie, _ = cookiejar.New(nil)

	client := &http.Client{
		Timeout: time.Second * DefaultTimeout,
		Jar:     cookie,
	}

	debug()
	return &Client{client}
}

// NewNoCookie return a new Client that won't save cookies
func NewNoCookie() *Client {
	client := &http.Client{
		Timeout: time.Second * DefaultTimeout,
	}

	debug()
	return &Client{client}
}

// Data is the body of http request
type Data map[string][]string

// Header is the header of http request
type Header map[string]string

// Option holds all the #Request requirements
type Option struct {
	URL      string // required
	Method   string // default: "GET", anything "POST", "PUT", "DELETE" or "PATCH"
	BodyStr  string
	Body     *Data
	Form     *Data       // set Content-Type header as "application/x-www-form-urlencoded"
	JSON     interface{} // set Content-Type header as "application/json"
	Query    *Data
	QueryRaw string
	Header   *Header
}

// SetTimeout sets client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	debug(timeout)

	c.httpClient.Timeout = timeout
}

// Request sends http request
func (c *Client) Request(opt *Option) (data []byte, res *http.Response, err error) {
	//set GET as default method
	if opt.Method == "" {
		opt.Method = "GET"
	}

	opt.Method = strings.ToUpper(opt.Method)

	//url
	reqURL, err := makeURL(opt.URL, opt.Query)
	if err != nil {
		return
	}

	//body
	reqBody, err := makeBody(opt)
	if err != nil {
		return
	}

	req, err := http.NewRequest(opt.Method, reqURL.String()+opt.QueryRaw, strings.NewReader(reqBody))
	if err != nil {
		debug("ERR(req)", err)
		return
	}

	//header
	makeHeader(req, opt)

	debug(req.Method, "\t>", req.URL.String())
	now := time.Now()

	res, err = c.httpClient.Do(req)
	if err != nil {
		debug("ERR", "\t<", err, humanizeNano(time.Now().Sub(now)))
		return
	}
	defer res.Body.Close()

	debug(res.StatusCode, "\t<", res.Request.URL, humanizeNano(time.Now().Sub(now)))

	// read all
	// it's a good practice to read all data so golang http can reuse requests
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		debug("ERR(ReadAll)", err)
		return
	}

	return
}

func makeURL(urlStr string, query *Data) (u *url.URL, err error) {
	u, err = url.Parse(urlStr)
	if err != nil {
		debug("ERR:", err)
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

func makeBody(opt *Option) (body string, err error) {
	var data *Data

	switch {
	case opt.BodyStr != "":
		body = opt.BodyStr
		return

	case opt.JSON != nil:
		jsonStr, err := json.Marshal(opt.JSON)
		if err != nil {
			debug("ERR:", err)
			return body, err
		}

		body = string(jsonStr)
		return body, err

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

func makeHeader(req *http.Request, opt *Option) {
	// default User-Agent
	// req.Header.Set("User-Agent", "github.com/ddo/request")
	req.Header.Set("User-Agent", " ") // == "" on the host side

	switch {
	case opt.Form != nil:
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	case opt.JSON != nil:
		req.Header.Set("Content-Type", "application/json")
	}

	if opt.Header == nil {
		return
	}

	for key, value := range *opt.Header {
		req.Header.Set(key, value)
	}
}

func humanizeNano(n time.Duration) string {
	var suffix string

	switch {
	case n > 1e9:
		n /= 1e9
		suffix = "s"
	case n > 1e6:
		n /= 1e6
		suffix = "ms"
	case n > 1e3:
		n /= 1e3
		suffix = "us"
	default:
		suffix = "ns"
	}

	return strconv.Itoa(int(n)) + suffix
}
