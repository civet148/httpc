package mock

import (
	"encoding/json"
	"fmt"
)

type AuthType string

const (
	AuthType_Null   AuthType = ""
	AuthType_Basic  AuthType = "Basic"
	AuthType_Bearer AuthType = "Bearer"
)

func (t AuthType) String() string {
	return string(t)
}
func (t AuthType) Valid() bool {
	switch t {
	case AuthType_Basic:
		return true
	case AuthType_Bearer:
		return true
	}
	return false
}

const (
	HEADER_AUTHORIZATION = "Authorization"
	HEADER_AUTH_TOKEN    = "Auth-Token"
)

const (
	MAX_DURATION int32  = 1440
	MIN_DURATION int32  = 1
	JSON_RPC_VER string = "2.0"
)

type BizCode int

const (
	CODE_ERROR                           BizCode = -1   //unknown error
	CODE_OK                              BizCode = 0    //success
	CODE_TOO_MANAY_REQUESTS              BizCode = 429  //too many requests
	CODE_INTERNAL_SERVER_ERROR           BizCode = 500  //internal service error
	CODE_DATABASE_ERROR                  BizCode = 501  //database server error
	CODE_ACCESS_DENY                     BizCode = 1000 //access deny
	CODE_UNAUTHORIZED                    BizCode = 1001 //user unauthorized
	CODE_INVALID_USER_OR_PASSWORD        BizCode = 1002 //user or password incorrect
	CODE_INVALID_PARAMS                  BizCode = 1003 //parameters invalid
	CODE_INVALID_JSON_OR_REQUIRED_PARAMS BizCode = 1004 //json format is invalid
	CODE_ALREADY_EXIST                   BizCode = 1005 //account name already exist
	CODE_NOT_FOUND                       BizCode = 1006 //record not found
	CODE_INVALID_PASSWORD                BizCode = 1007 //wrong password
	CODE_INVALID_AUTH_CODE               BizCode = 1008 //invalid auth code
	CODE_ACCESS_VIOLATE                  BizCode = 1009 //access violate
	CODE_TYPE_UNDEFINED                  BizCode = 1010 //type undefined
	CODE_BAD_DID_OR_SIGNATURE            BizCode = 1011 //bad did or signature
	CODE_ACCOUNT_BANNED                  BizCode = 1012 //account was banned
	CODE_EXPORT_FAILED                   BizCode = 1013
	CODE_PRIVILEGE_EXIT                  BizCode = 1014
	CODE_USERNAME_EXIT                   BizCode = 1015
	CODE_PHONENUMBER_EXIT                BizCode = 1016
	CODE_PHONE_FORMAT_WRONG              BizCode = 1017
	CODE_ROLE_EXIT                       BizCode = 1018
	CODE_SESSION_CONTEXT                 BizCode = 1020
	CODE_NO_PRIVILEGE                    BizCode = 1021
	CODE_GET_PRIVILEGE_FAILED            BizCode = 1022
	CODE_CHAIN_VERFICATION_FAILED        BizCode = 1023
	CODE_EDIT_WINDING_CYCLE_FAILED       BizCode = 1024
	CODE_EDIT_CONTROL_FAILED             BizCode = 1025
	CODE_MAX_DURATION_LIMIT              BizCode = 1026
	CODE_EMAIL_FORMAT_WRONG              BizCode = 1027
	CODE_EMAIL_EXIT                      BizCode = 1028
	CODE_EMAIL_SEND_FAILED               BizCode = 1029
	CODE_EMAIL_CODE_TIMEOUT              BizCode = 1030
	CODE_EMAIL_NOT_EXIT                  BizCode = 1031
	CODE_EMAIL_CODE_ERROR                BizCode = 1032
	CODE_USER_INACTIVE                   BizCode = 1033
	CODE_CHAIN_ACCREDIT_ERROR            BizCode = 1034
	CODE_DEFAULT_POOL_ALREADY_EXIST      BizCode = 1035
)

var codeMessages = map[BizCode]string{
	CODE_ERROR:                           "ERROR",
	CODE_OK:                              "OK",
	CODE_TOO_MANAY_REQUESTS:              "请求频率超限",
	CODE_INTERNAL_SERVER_ERROR:           "内部服务器错误",
	CODE_DATABASE_ERROR:                  "数据库错误",
	CODE_UNAUTHORIZED:                    "未认证授权",
	CODE_ACCESS_DENY:                     "访问拒绝",
	CODE_INVALID_USER_OR_PASSWORD:        "用户名/密码错误",
	CODE_INVALID_PARAMS:                  "无效参数",
	CODE_INVALID_JSON_OR_REQUIRED_PARAMS: "无效的JSON请求格式",
	CODE_ALREADY_EXIST:                   "数据已存在，不能重复添加",
	CODE_NOT_FOUND:                       "数据没找到",
	CODE_INVALID_PASSWORD:                "无效密码",
	CODE_INVALID_AUTH_CODE:               "无效验证码",
	CODE_ACCESS_VIOLATE:                  "访问违规",
	CODE_TYPE_UNDEFINED:                  "未定义",
	CODE_BAD_DID_OR_SIGNATURE:            "无效的数字身份/签名",
	CODE_ACCOUNT_BANNED:                  "账号已禁用",
	CODE_EXPORT_FAILED:                   "导出失败，请检查网络情况/联系管理员",
	CODE_PRIVILEGE_EXIT:                  "该用户已有权限，请勿重复添加",
	CODE_USERNAME_EXIT:                   "该用户帐号已存在，请您更换用户帐号",
	CODE_PHONENUMBER_EXIT:                "该手机号码已存在，请您更换号码",
	CODE_PHONE_FORMAT_WRONG:              "手机号码格式错误",
	CODE_ROLE_EXIT:                       "该角色名称已存在",
	CODE_SESSION_CONTEXT:                 "会话超时",
	CODE_NO_PRIVILEGE:                    "无权限",
	CODE_GET_PRIVILEGE_FAILED:            "查询权限失败",
	CODE_CHAIN_VERFICATION_FAILED:        "数据校验失败",
	CODE_EDIT_WINDING_CYCLE_FAILED:       "更新上链周期失败",
	CODE_EDIT_CONTROL_FAILED:             "编辑禁用/启用失败",
	CODE_MAX_DURATION_LIMIT:              fmt.Sprintf("上链周期区间为%d-%d分钟 请您输入合理区间内", MIN_DURATION, MAX_DURATION),
	CODE_EMAIL_FORMAT_WRONG:              "邮箱格式错误",
	CODE_EMAIL_EXIT:                      "该邮箱已存在，请您更换邮箱",
	CODE_EMAIL_SEND_FAILED:               "邮件发送失败！",
	CODE_EMAIL_CODE_TIMEOUT:              "验证码超时！",
	CODE_EMAIL_NOT_EXIT:                  "邮箱不存在",
	CODE_EMAIL_CODE_ERROR:                "验证码错误！",
	CODE_USER_INACTIVE:                   "用户未激活或已禁用",
	CODE_CHAIN_ACCREDIT_ERROR:            "区块链数据交互失败",
	CODE_DEFAULT_POOL_ALREADY_EXIST:      "网络对应默认节点池已存在",
}

// 错误提示
const (
	ERROR_MESSAGE_SERVER_ERROR       = "Server Error"
	ERROR_MESSAGE_SESSION_ERROR      = "User session context"
	ERROR_MESSAGE_NO_PROVILEGE       = "No privilege"
	ERROR_MESSAGE_PARAMETER_ERROE    = "Parameter error"
	ERROR_MESSAGE_LOGIN_FAILED       = "UserName or password error"
	ERROR_MESSAGE_NO_ROLE            = "User hadn't role"
	ERROR_MESSAGE_LIST_QUERY_FAILED  = "List query failed"
	ERROR_MESSAGE_USER_EXIT          = "User already exists"
	ERROR_MESSAGE_PHONE_EXIT         = "Phone already exists"
	ERROR_MESSAGE_EMAIL_EXIT         = "Email already exists"
	ERROR_MESSAGE_AUTH_FAILED        = "Authorization failed"
	ERROR_MESSAGE_GET_AUTH_FAILED    = "Get privilege failed"
	ERROR_MESSAGE_EXPORT_FAILED      = "Data export filed"
	ERROR_MESSAGE_CHECK_CHAIN_FAILED = "Prefix verification failed"
	ERROR_MESSAGE_EDIT_CYCLE_FAILED  = "Edit winding cycle failed"
	ERROR_MESSAGE_CONTROL_FAILE      = "Enable/Disable Failed"
	ERROR_MESSAGE_BAD_DID_OR_SIGN    = "Bad did or signature"
	ERROR_MESSAGE_ACCOUNT_BANNED     = "Account was banned"
)

func (c BizCode) Ok() bool {
	return c == CODE_OK
}

func (c BizCode) String() string {
	if m, ok := codeMessages[c]; ok {
		return m
	}
	return fmt.Sprintf("CODE_UNKNOWN<%d>", c)
}

func (c BizCode) GoString() string {
	return c.String()
}

type HttpHeader struct {
	Code    BizCode `json:"code"`    //response code of business (0=OK, other fail)
	Message string  `json:"message"` //error message
	Total   int64   `json:"total"`   //result total
	Count   int     `json:"count"`   //result count (single page)
}

type HttpResponse struct {
	Header HttpHeader  `json:"header"` //response header
	Data   interface{} `json:"data"`   //response data body
}

type RpcRequest struct {
	Id      interface{} `json:"id"`       //0
	JsonRpc string      `json:"json_rpc"` //2.0
	Method  string      `json:"method"`   //JSON-RPC method
	//Params  []interface{} `json:"params"`   //JSON-RPC parameters [any...]
}

type RpcError struct {
	Code    BizCode     `json:"code"`    //response code of business (0=OK, other fail)
	Message string      `json:"message"` //error message
	Data    interface{} `json:"data"`    //error attach data
}

type RpcResponse struct {
	Id      interface{} `json:"id"`       //0
	JsonRpc string      `json:"json_rpc"` //2.0
	Error   RpcError    `json:"error"`    //error message
	Result  interface{} `json:"result"`   //JSON-RPC result
}

func (r *RpcResponse) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

func NewRpcResponse(id interface{}, result interface{}) *RpcResponse {
	return &RpcResponse{
		Id:      id,
		JsonRpc: JSON_RPC_VER,
		Result:  result,
	}
}

func NewRpcError(id interface{}, code BizCode, strError string) *RpcResponse {
	return &RpcResponse{
		Id:      id,
		JsonRpc: JSON_RPC_VER,
		Error: RpcError{
			Code:    code,
			Message: strError,
			Data:    nil,
		},
	}
}
