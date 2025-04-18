package main

import (
	"encoding/json"
	"github.com/civet148/httpc"
	"github.com/civet148/log"
	"github.com/valyala/fastjson"
)

func init() {
	log.SetLevel("debug")
}

func main() {
	c := &httpc.Client{}
	defer c.Close()
	//FastJson()
	FilfoxGet(c)
}

func FilfoxGet(c *httpc.Client) {
	c.Header().Set("token", "12345678901234567890")

	uv := httpc.NewUrlValues()
	uv.Add("page", 0).Add("pageSize", 15)

	r, err := c.Get("https://filfox.info/api/v1/address/f07749/blocks", uv.Values())
	log.Debugf("Status code [%v] content type [%s] data [%+v]", r.StatusCode, r.ContentType, string(r.Body))
	if err != nil {
		log.Errorf("GET error [%s]", err)
		return
	}
	//
	//var body = make(url.Values)
	//body.Add("block_id", "9566631935999642852")
	//var query = make(url.Values)
	//query.Add("sign_type", "R S A2")
	//query.Add("version", "1.0")
	//r, err = c.PostFormUrlEncoded("https://filfox.info/api/v1/address/f07749/blocks", body, query)
	//if err != nil {
	//	log.Errorf("POST error [%s]", err)
	//	return
	//}
	//log.Debugf("POST response [%s]", r.Body)
}

func FastJson() {
	type RespHeader struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Count   int    `json:"count"`
		Total   int    `json:"total"`
	}
	type RespData struct {
		Name string `json:"name"`
		Sex  string `json:"sex"`
		Age  int    `json:"age"`
	}

	strJson := `
{
   "header":{
      "code":0,
      "message":"OK",
      "count":1,
      "total":1
   },
   "data":[{
 		"name":"lory",
        "sex":"male",
        "age":18
   },{
 		"name":"jhon",
        "sex":"male",
        "age":28
   }]
}
`
	values, err := fastjson.Parse(strJson)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	var data []RespData
	v := values.Get("data")
	if v != nil {
		var body []byte
		body = v.MarshalTo(nil)
		log.Infof("body [%s]", body)
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		log.Infof("response data %+v", data)
	}
}

//github.com/valyala/fastjson
