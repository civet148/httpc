package mock

import (
	"encoding/json"
	"github.com/civet148/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

type Controller struct {
	cfg *Config
}

func NewController(cfg *Config) *Controller {
	return &Controller{
		cfg: cfg,
	}
}

func (m *Controller) OK(c *gin.Context, data interface{}, count int, total int64) {
	var code = CODE_OK
	if data == nil {
		data = struct{}{}
	}
	var r = &HttpResponse{
		Header: HttpHeader{
			Code:    code,
			Message: code.String(),
			Count:   count,
			Total:   total,
		},
		Data: data,
	}
	if data != nil {
		log.Json("response", r)
	}
	c.JSON(http.StatusOK, r)
	c.Abort()
}

func (m *Controller) Error(c *gin.Context, code BizCode, message string) {
	var r = &HttpResponse{
		Header: HttpHeader{
			Code:    code,
			Message: message,
			Count:   0,
		},
		Data: struct{}{},
	}
	log.Errorf("[Controller] response error code [%d] message [%s]", code, message)
	c.JSON(http.StatusOK, r)
	c.Abort()
}

func (m *Controller) ErrorStatus(c *gin.Context, status int, message string) {
	log.Errorf("[Controller] http status code [%d] message [%s]", status, message)
	c.String(status, message)
	c.Abort()
}

func (m *Controller) RpcResult(c *gin.Context, data interface{}, err error, id interface{}) {
	var status = http.StatusOK
	var strResp string
	if err != nil {
		status = http.StatusInternalServerError
		data = &RpcResponse{
			Id:      id,
			JsonRpc: "2.0",
			Error: RpcError{
				Code:    500,
				Message: err.Error(),
			},
			Result: nil,
		}
	}
	switch data.(type) {
	case string:
		strResp = data.(string)
	default:
		{
			b, _ := json.Marshal(data)
			strResp = string(b)
		}
	}

	c.String(status, strResp)
	c.Abort()
}

func (m *Controller) GetClientIP(c *gin.Context) (strIP string) {
	return c.ClientIP()
}

func (m *Controller) Authorization(c *gin.Context) (authType AuthType, strKey, strSecret string) {
	strAuth := c.Request.Header.Get(HEADER_AUTHORIZATION)
	if strAuth == "" {
		log.Warnf("no header -> authorization found")
		return AuthType_Null, "", ""
	}
	var ok bool
	strKey, strSecret, ok = c.Request.BasicAuth()
	if ok {
		authType = AuthType_Basic
	} else {
		authType = AuthType_Bearer
		strToken := strings.TrimPrefix(strAuth, AuthType_Bearer.String()+" ")
		ss := strings.Split(strToken, ".")
		count := len(ss)
		if count <= 1 {
			strKey = ss[0]
		} else if count == 2 {
			strKey = ss[0]
			strSecret = ss[1]
		} else {
			strKey = ss[1]
			strSecret = ss[2]
		}
	}
	log.Debugf("ip [%s] auth type [%s] key [%s] secret [%s] auth [%s]", c.ClientIP(), authType, strKey, strSecret, strAuth)
	return
}

func (m *Controller) bindJSON(c *gin.Context, req interface{}) (err error) {
	if err = c.ShouldBindJSON(req); err != nil {
		log.Errorf("%s", err)
		m.Error(c, CODE_INVALID_JSON_OR_REQUIRED_PARAMS, ERROR_MESSAGE_PARAMETER_ERROE)
		c.Abort()
		return
	}

	body, _ := json.MarshalIndent(req, "", "\t")
	log.Debugf("request from [%s] body [%+v]", c.ClientIP(), string(body))
	return
}

func (m *Controller) WebSocketRpcV1(c *gin.Context) {
	uri := c.Request.RequestURI
	_, strKey, strSecret := m.Authorization(c)
	log.Debugf("uri [%s] websocket client ip [%s] auth key [%s] secret [%s]", c.ClientIP(), uri, strKey, strSecret)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("web socket upgrade error [%s]", err)
		return
	}
	ws := NewWebSocket(conn)
	defer ws.Close()
	for {
		var mt int
		var message []byte

		mt, message, err = ws.ReadMessage()
		if err != nil {
			log.Warnf("websocket client is closed")
			break
		}

		var msg *RpcMessage
		msg, err = MakeRpcMessage(uri, message)
		if err != nil {
			log.Errorf("make rpc message error [%s]", err)
			break
		}

		log.Debugf("websocket session id [%s] method [%s]", msg.SessionId, msg.RpcMethod)
		err = m.websocketRpcCall(ws, mt, msg)
		if err != nil {
			log.Errorf(err.Error())
			break
		}
	}
}

func (m *Controller) websocketRpcCall(ws *WebSocket, msgType int, msg *RpcMessage) (err error) {
	var strResp string
	strResp = m.cfg.DataRaw
	log.Debugf("websocket request msg [%+v]", msg)
	err = ws.WriteMessage(msgType, []byte(strResp))
	if err != nil {
		return log.Errorf("write message error [%s]", err)
	}
	log.Debugf("websocket response [%+v]", strResp)
	return nil
}
