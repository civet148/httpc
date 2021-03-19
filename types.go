package httpc

import "net/http"

const (
	HTTP_METHOD_GET     = http.MethodGet     //请求指定的页面信息，并返回实体主体。
	HTTP_METHOD_POST    = http.MethodPost    //向指定资源提交数据进行处理请求（例如提交表单或者上传文件）。数据被包含在请求体中。POST 请求可能会导致新的资源的建立和/或已有资源的修改。
	HTTP_METHOD_PUT     = http.MethodPut     //从客户端向服务器传送的数据取代指定的文档的内容。
	HTTP_METHOD_CONNECT = http.MethodConnect //HTTP/1.1 协议中预留给能够将连接改为管道方式的代理服务器
	HTTP_METHOD_OPTIONS = http.MethodOptions //允许客户端查看服务器的性能。
	HTTP_METHOD_DELETE  = http.MethodDelete  //请求服务器删除指定的页面内容。
	HTTP_METHOD_TRACE   = http.MethodTrace   //回显服务器收到的请求，主要用于测试或诊断。
	HTTP_METHOD_PATCH   = http.MethodPatch   //是对PUT方法的补充，用来对已知资源进行局部更新。
	HTTP_METHOD_HEAD    = http.MethodHead    //类似于GET请求，只不过返回的响应中没有具体的内容，用于获取报头
)

const (
	HEADER_KEY_CONTENT_TYPE  = "Content-Type"
	HEADER_KEY_AUTHORIZATION = "Authorization"
)

const (
	CONTENT_TYPE_NAME_TEXT_PLAIN             = "text/plain"                        //content-type (raw)
	CONTENT_TYPE_NAME_MULTIPART_FORM_DATA    = "multipart/form-data"               //content-type (form-data)
	CONTENT_TYPE_NAME_X_WWW_FORM_URL_ENCODED = "application/x-www-form-urlencoded" //content-type (urlencoded)
	CONTENT_TYPE_NAME_APPLICATION_JSON       = "application/json"                  //content-type (json)
	CONTENT_TYPE_NAME_TEXT_HTML              = "text/html"                         //content-type (html)
)

type UrlValues map[string]interface{} //POST提交表单参数
