package mock

import (
	"encoding/json"
	"github.com/civet148/gotools/randoms"
	"github.com/civet148/log"
)

type RpcMessage struct {
	SessionId     string      `json:"session_id"`
	Uri           string      `json:"uri"`
	RpcId         interface{} `json:"rpc_id"`
	RpcMethod     string      `json:"rpc_method"`
	RpcData       []byte      `json:"rpc_data"`
	PoolId        int32       `json:"pool_id"`
	ProjectId     int32       `json:"project_id"`
	ProjectName   string      `json:"project_name"`
	ProjectKey    string      `json:"project_key"`
	ProjectSecret string      `json:"project_secret"`
}

func MakeSessionId() string {
	return randoms.RandomAlphaOrNumeric(8, true, true)
}


func MakeRpcMessage(strUri string, message []byte) (msg *RpcMessage, err error) {
	var req RpcRequest
	if err = json.Unmarshal(message, &req); err != nil {
		log.Errorf(err.Error())
		return
	}
	msg = &RpcMessage{
		Uri:           strUri,
		RpcId:         req.Id,
		RpcMethod:     req.Method,
		RpcData:       message,
		SessionId:      MakeSessionId(),
	}
	return msg, nil
}
