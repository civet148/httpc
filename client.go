package httpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/civet148/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	cli    *http.Client
	header *Header
}

type Response struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

func init() {
	log.SetLevel(log.LEVEL_INFO)
}

//new a normal http client with timeout (seconds)
func NewHttpClient(timeout int) *Client {
	return newClient(timeout)
}

//new a https client
//timeout   timeout seconds
//cer       PEM certification path
func NewHttpsClient(timeout int, cer interface{}) *Client {
	return newClient(timeout, cer)
}

func newClient(timeout int, args ...interface{}) (c *Client) {

	if timeout <= 0 {
		timeout = 3
	}
	c = &Client{
		header: &Header{
			values: map[string]string{
				HEADER_KEY_CONTENT_TYPE: CONTENT_TYPE_NAME_X_WWW_FORM_URL_ENCODED,
			},
		},
		cli: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
	return
}

func (c *Client) Debug() {
	log.SetLevel(0)
}

func (c *Client) Header() *Header {
	return c.header
}

//send a http request by GET method
func (c *Client) Get(strUrl string, values UrlValues) (response *Response, err error) {
	return c.do(HTTP_METHOD_GET, strUrl, values)
}

//send a http request by POST method with application/x-www-form-urlencoded
func (c *Client) PostUrlEncoded(strUrl string, values UrlValues) (response *Response, err error) {
	return c.do(HTTP_METHOD_POST, strUrl, values)
}

//send a http request by PUT method
func (c *Client) Put(strUrl string) (response *Response, err error) {
	return c.do(HTTP_METHOD_PUT, strUrl, nil)
}

//send a http request by DELETE method
func (c *Client) Delete(strUrl string) (response *Response, err error) {
	return c.do(HTTP_METHOD_DELETE, strUrl, nil)
}

//send a http request by TRACE method
func (c *Client) Trace(strUrl string) (response *Response, err error) {
	return c.do(HTTP_METHOD_TRACE, strUrl, nil)
}

//send a http request by PATCH method
func (c *Client) Patch(strUrl string) (response *Response, err error) {
	return c.do(HTTP_METHOD_PATCH, strUrl, nil)
}

//send a http request by POST method with content-type specified
//data type must could be string,[]byte,UrlValues,struct and so on
func (c *Client) Post(strContentType string, strUrl string, data interface{}) (response *Response, err error) {
	c.header.Set(HEADER_KEY_CONTENT_TYPE, strContentType)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//send a http request by POST method with content-type application/json
//data type must could be string,[]byte,UrlValues,struct and so on
func (c *Client) PostJson(strUrl string, data interface{}) (response *Response, err error) {
	c.header.Set(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_APPLICATION_JSON)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//send a http request by POST method with content-type text/plain
//data type must could be string,[]byte,UrlValues,struct and so on
func (c *Client) PostRaw(strUrl string, data interface{}) (response *Response, err error) {
	c.header.Set(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_TEXT_PLAIN)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//send a http request by POST method with content-type multipart/form-data
//data type must could be string,[]byte,UrlValues,struct and so on
func (c *Client) PostFormData(strUrl string, data interface{}) (response *Response, err error) {
	c.header.Set(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_MULTIPART_FORM_DATA)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//send a http request by POST method with content-type multipart/form-data
//data type must could be string,[]byte,UrlValues,struct and so on
func (c *Client) PostFormUrlEncoded(strUrl string, data interface{}) (response *Response, err error) {
	c.header.Set(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_X_WWW_FORM_URL_ENCODED)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//do send request to destination host
func (c *Client) do(strMethod, strUrl string, data interface{}) (response *Response, err error) {

	var body io.Reader

	if data != nil {

		switch data.(type) {
		case UrlValues:
			{
				kvs := url.Values{}
				values := data.(UrlValues)
				if values != nil {
					for k, v := range values {
						kvs.Add(k, fmt.Sprintf("%v", v))
					}
				}
				body = strings.NewReader(kvs.Encode())
				log.Debug("UrlValues -> [%+v]", data.(UrlValues))
			}
		case string:
			{
				body = strings.NewReader(data.(string))
				log.Debug("string -> [%s]", data.(string))
			}
		case []byte:
			{
				body = bytes.NewReader(data.([]byte))
				log.Debug("[]byte -> [%s]", data.([]byte))
			}
		default:
			{
				var jsonData []byte
				if jsonData, err = json.Marshal(data); err != nil {
					log.Error("can't marshal data to json, error [%v]", err.Error())
					return
				}
				body = bytes.NewReader(jsonData)
				log.Debug("object -> [%s]", jsonData)
			}
		}
	}

	if response, err = c.sendRequest(strMethod, strUrl, body); err != nil {
		return
	}
	return
}

func (c *Client) sendRequest(strMethod, strUrl string, body io.Reader) (response *Response, err error) {

	var req *http.Request
	var resp *http.Response

	if req, err = http.NewRequest(strMethod, strUrl, body); err != nil {
		log.Errorf("new request error [%s]", err)
		return
	}

	for k, v := range c.header.values {
		req.Header.Set(k, v)
	}

	if resp, err = c.cli.Do(req); err != nil {
		log.Errorf("send request error [%s]", err)
		return
	}

	defer resp.Body.Close()

	response = &Response{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get(HEADER_KEY_CONTENT_TYPE),
	}

	if response.Body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Errorf("%s", err)
		return
	}
	return
}
