package httpc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/civet148/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type uploadFile struct {
	Name     string
	FilePath string
}

type Client struct {
	cli    http.Client
	header http.Header
	locker sync.RWMutex
}

func init() {
	log.SetLevel(log.LEVEL_INFO)
}

//new a normal http client with timeout (seconds)
func NewHttpClient(opts ...*Option) *Client {

	return newClient(opts...)
}

func newClient(opts ...*Option) (c *Client) {
	var header http.Header
	var tlsConf *tls.Config
	var opt = &Option{
		Timeout: 30,
	}
	for _, o := range opts {
		opt = o
	}
	if opt != nil {
		header = opt.Header
		tlsConf = opt.TlsConf
	}
	if header == nil {
		header = http.Header{}
	}
	if tlsConf == nil {
		tlsConf = &tls.Config{InsecureSkipVerify: true}
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConf,
	}
	c = &Client{
		header: header,
		cli: http.Client{
			Transport: tr,
			Timeout:   time.Duration(opt.Timeout) * time.Second,
		},
	}
	return
}

func (c *Client) Close() {
	c.cli.CloseIdleConnections()
}

func (c *Client) Debug() {
	log.SetLevel(0)
}

func (c *Client) Header() http.Header {
	return c.header
}

func (c *Client) SetBasicAuth(username, password string) {
	c.setHeader("Authorization", "Basic "+basicAuth(username, password))
}

func (c *Client) SetOAuth2(token string) {
	c.SetBearerToken(token)
}

func (c *Client) SetBearerToken(token string) {
	c.setHeader("Authorization", "Bearer "+token)
}

func (c *Client) GetEx(strUrl string, values url.Values, v interface{}) (status int, err error) {
	var r *Response
	if r, err = c.Get(strUrl, values); err != nil {
		return http.StatusBadGateway, err
	}
	if r.StatusCode == http.StatusOK {
		if err = r.Unmarshal(v); err != nil {
			log.Errorf("unmarshal response body data to struct error [%s]", err)
			return r.StatusCode, err
		}
	}
	return r.StatusCode, nil
}

//send a http request by GET method
func (c *Client) Get(strUrl string, values url.Values) (r *Response, err error) {
	return c.get(strUrl, values)
}

//send a http request by GET method and copy to writter
func (c *Client) CopyFile(strUrl string, writer io.Writer) (written int64, err error) {
	var r *http.Response
	r, err = http.Get(strUrl)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()
	written, err = io.Copy(writer, r.Body)
	if err != nil {
		return 0, err
	}
	return
}

//send a http request by GET method and save to file
func (c *Client) SaveFile(strUrl string, strFilePath string) (written int64, err error) {
	var r *http.Response
	r, err = http.Get(strUrl)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()
	var dst *os.File
	dst, err = os.Create(strFilePath)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	written, err = io.Copy(dst, r.Body)
	if err != nil {
		return 0, err
	}
	return
}

//send a http request by POST method with application/x-www-form-urlencoded
func (c *Client) PostUrlEncoded(strUrl string, values url.Values) (r *Response, err error) {
	return c.do(HTTP_METHOD_POST, strUrl, values)
}

//send a http request by GET method and unmarshal json data to struct v
func (c *Client) GetJson(strUrl string, values url.Values, v interface{}) (status int, err error) {
	var r *Response
	if r, err = c.get(strUrl, values); err != nil {
		log.Errorf("GET url [%s] values [%+v] error [%s]", strUrl, values, err.Error())
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET url [%s] values [%+v] remote server status code [%v]", strUrl, values, r.StatusCode)
		log.Errorf(err.Error())
		return r.StatusCode, err
	}
	//log.Debugf("url [%s] values [%+v] response [%s]", strUrl, values, string(r.Body))
	if err = json.Unmarshal(r.Body, v); err != nil {
		log.Errorf("json unmarshal error [%s] data body [%s]", err, r.Body)
		return
	}
	return http.StatusOK, nil
}

//send a http request by PUT method
func (c *Client) Put(strUrl string, body io.Reader) (r *Response, err error) {
	return c.SendRequest(c.header, HTTP_METHOD_PUT, strUrl, body)
}

//send a http request by DELETE method
func (c *Client) Delete(strUrl string) (r *Response, err error) {
	return c.do(HTTP_METHOD_DELETE, strUrl, nil)
}

//send a http request by TRACE method
func (c *Client) Trace(strUrl string) (r *Response, err error) {
	return c.do(HTTP_METHOD_TRACE, strUrl, nil)
}

//send a http request by PATCH method
func (c *Client) Patch(strUrl string) (r *Response, err error) {
	return c.do(HTTP_METHOD_PATCH, strUrl, nil)
}

//send a http request by POST method with content-type specified
//data type could be string,[]byte,url.Values,struct and so on
func (c *Client) Post(strContentType string, strUrl string, data interface{}) (r *Response, err error) {
	c.setContentType(strContentType)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//send a http request by POST method with content-type application/json
//data type could be string,[]byte,url.Values,struct and so on and
func (c *Client) PostJson(strUrl string, data interface{}) (r *Response, err error) {
	c.setContentType(CONTENT_TYPE_NAME_APPLICATION_JSON)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//send a http request by POST method with content-type text/plain
//data type must could be string,[]byte,url.Values,struct and so on
func (c *Client) PostRaw(strUrl string, data interface{}) (r *Response, err error) {
	c.setContentType(CONTENT_TYPE_NAME_TEXT_PLAIN)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

//send a http request by POST method with content-type multipart/form-data
//data type must could be string,[]byte,url.Values,struct and so on
func (c *Client) PostFormData(strUrl string, data interface{}) (r *Response, err error) {
	c.setContentType(CONTENT_TYPE_NAME_MULTIPART_FORM_DATA)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

/*
send a http request by POST method with content-type multipart/form-data
kvs a map of key=value, if the value is a file path please use @ as prefix
example:

var params = map[string]string{
      "image_name":"a.jpg",
      "image_file":"@/tmp/a.jpg"
}
*/
func (c *Client) PostFormDataMultipart(strUrl string, params map[string]string) (r *Response, err error) {
	c.setContentType(CONTENT_TYPE_NAME_MULTIPART_FORM_DATA)
	return c.doUpload(strUrl, params)
}

//send a http request by POST method with content-type multipart/form-data
//data type must could be string,[]byte,url.Values,struct and so on
func (c *Client) PostFormUrlEncoded(strUrl string, data interface{}) (r *Response, err error) {
	c.setContentType(CONTENT_TYPE_NAME_X_WWW_FORM_URL_ENCODED)
	return c.do(HTTP_METHOD_POST, strUrl, data)
}

func (c *Client) setHeader(key, value string) {
	c.locker.Lock()
	if c.header == nil {
		c.header = http.Header{}
	}
	c.header.Set(key, value)
	c.locker.Unlock()
}

func (c *Client) setContentType(contentType string) {
	c.setHeader(HEADER_KEY_CONTENT_TYPE, contentType)
}

//do send request to destination host
func (c *Client) do(strMethod, strUrl string, data interface{}) (r *Response, err error) {

	var body io.Reader

	if data != nil {

		switch data.(type) {
		case url.Values:
			{
				values := data.(url.Values)
				body = strings.NewReader(values.Encode())
				log.Debug("url.Values -> [%+v]", values)
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

	if r, err = c.SendRequest(c.header, strMethod, strUrl, body); err != nil {
		return
	}
	return
}

func (c *Client) get(strUrl string, values url.Values) (r *Response, err error) {

	if values != nil {
		u, err := url.Parse(strUrl)
		if err != nil {
			return nil, err
		}
		u.RawQuery = values.Encode()
		strUrl = u.String()
	}
	return c.SendRequest(c.header, HTTP_METHOD_GET, strUrl, nil)
}

func (c *Client) SendRequest(header http.Header, strMethod, strUrl string, body io.Reader) (r *Response, err error) {

	var req *http.Request
	var resp *http.Response

	if req, err = http.NewRequest(strMethod, strUrl, body); err != nil {
		log.Errorf("new request error [%s]", err)
		return
	}

	req.Header = header

	if resp, err = c.cli.Do(req); err != nil {
		log.Errorf("send request error [%s]", err)
		return
	}

	defer resp.Body.Close()

	r = &Response{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get(HEADER_KEY_CONTENT_TYPE),
	}

	if r.Body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Errorf("%s", err)
		return
	}
	return
}

func (c *Client) doUpload(strUrl string, params map[string]string) (r *Response, err error) {
	var body io.Reader
	var contentType string
	body, contentType, err = c.getReader(params)
	if err != nil {
		return nil, log.Errorf(err.Error())
	}
	c.setHeader(HEADER_KEY_CONTENT_TYPE, contentType)
	return c.SendRequest(c.header, HTTP_METHOD_POST, strUrl, body)
}

func (c *Client) getReader(params map[string]string) (reader io.Reader, contentType string, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	for k, v := range params {
		if strings.HasPrefix(v, "@") {
			path := v[1:]
			upfile := &uploadFile{
				Name:     k,
				FilePath: path,
			}
			var file *os.File

			file, err = os.Open(upfile.FilePath)
			if err != nil {
				return reader, "", log.Errorf("open file [%s] error [%s]", upfile.FilePath, err)
			}
			defer file.Close()
			log.Debugf("open file [%s] ok", upfile.FilePath)
			var part io.Writer
			part, err = writer.CreateFormFile(upfile.Name, filepath.Base(upfile.FilePath))
			log.Debugf("writer.CreateFormFile field name [%s] file name [%s]", upfile.Name, filepath.Base(upfile.FilePath))
			if err != nil {
				return reader, "", log.Errorf("writer.CreateFormFile field name [%s] file name [%s] error [%s]", upfile.Name, filepath.Base(upfile.FilePath), err)
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return reader, "", log.Errorf("io.Copy error [%s]", err)
			}
		} else {
			err = writer.WriteField(k, v)
			if err != nil {
				return reader, "", log.Errorf("write key %s value %s error %s", k, v, err.Error())
			}
		}
	}

	return body, writer.FormDataContentType(), nil
}
