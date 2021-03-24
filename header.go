package httpc

import (
	"fmt"
	"sync"
)

type Header struct {
	mutex  sync.Mutex
	values map[string]string
}

func (h *Header) find(strKey string) (strValue string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.values[strKey]
}

func (h *Header) add(strKey, strValue string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.values[strKey] = strValue
}

func (h *Header) del(strKey string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.values, strKey)
}

func (h *Header) clear() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	for k, _ := range h.values {
		delete(h.values, k)
	}
}

//删除http请求头部参数
func (h *Header) Delete(strKey string) {
	h.del(strKey)
}

//删除http请求所有头部参数
func (h *Header) RemoveAll() {
	h.clear()
}

//获取http请求头部strKey对应的值
func (h *Header) Get(strKey string) (strValue string) {

	return h.find(strKey)
}

//设置http请求头部内容类型(自定义key-value)
func (h *Header) Set(strKey, strValue string) *Header {
	h.add(strKey, strValue)
	return h
}

//设置http请求头部内容类型(Content-Type=application/x-www-form-urlencoded)
func (h *Header) SetUrlEncoded() *Header {
	h.add(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_X_WWW_FORM_URL_ENCODED)
	return h
}

//设置http请求头部内容类型(Content-Type=multipart/form-data)
func (h *Header) SetMultipartFormData() *Header {
	h.add(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_MULTIPART_FORM_DATA)
	return h
}

//设置http请求头部内容类型(Content-Type=application/json)
func (h *Header) SetApplicationJson() *Header {
	h.add(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_APPLICATION_JSON)
	return h
}

//设置http请求头部内容类型(Content-Type=text/plain)
func (h *Header) SetTextPlain() *Header {
	h.add(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_TEXT_PLAIN)
	return h
}

//设置http请求头部内容类型(Content-Type=text/html)
func (h *Header) SetTextHtml() *Header {
	h.add(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_TEXT_HTML)
	return h
}

//设置http请求头部内容类型(Content-Type=application/octet-stream)
func (h *Header) SetOctetStream() *Header {
	h.add(HEADER_KEY_CONTENT_TYPE, CONTENT_TYPE_NAME_OCTET_STREAM)
	return h
}

//设置http请求头部key="Authorization" value=strValue
func (h *Header) SetAuthorization(strValue string) *Header {
	h.add(HEADER_KEY_AUTHORIZATION, strValue)
	return h
}

//设置http请求头部key="Content-Length" value=strValue
func (h *Header) SetContentLength(size int) *Header {
	h.add(HEADER_KEY_CONTENT_LENGTH, fmt.Sprintf("%v", size))
	return h
}
