package main

import (
	"encoding/json"
	"github.com/civet148/httpc"
	"github.com/civet148/log"
	"github.com/valyala/fastjson"
	"net/url"
	"time"
)

func init() {
	log.SetLevel("debug")
}

type Params struct {
	PageNo   int `json:"page_no"`
	PageSize int `json:"page_size"`
}

func main() {

	c := &httpc.Client{}
	defer c.Close()

	//FastJson()
	FilfoxGet(c)
}

func FilfoxGet(c *httpc.Client) {
	r, err := c.Get("https://filfox.info/api/v1/address/f07749/blocks", url.Values{
		"page":     []string{"0"},
		"pageSize": []string{"5"},
	})
	log.Debugf("Status code [%v] content type [%s] data [%+v]", r.StatusCode, r.ContentType, string(r.Body))
	if err != nil {
		log.Errorf("GET error [%s]", err)
		return
	}
	for i := 0; i < 5; i++ {
		r, err = c.Get("https://filfox.info/api/v1/address/f07749/blocks?page=1&pageSize=5", nil)
		if err != nil {
			log.Errorf("GET error [%s]", err)
			return
		}
		log.Debugf("[%d] response code [%v] content type [%s] data [%+v]", i, r.StatusCode, r.ContentType, string(r.Body))
		time.Sleep(2 * time.Second)
	}
	values := httpc.MakeQueryParams(&Params{
		PageNo:   1,
		PageSize: 100,
	})

	log.Infof("%+v", values)
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
